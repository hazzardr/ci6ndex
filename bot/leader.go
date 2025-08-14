package bot

import (
	"bytes"
	"ci6ndex/ci6ndex/generated"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	md "github.com/nao1215/markdown"
)

var (
	// 1 based page number. equivalent to offset+1
	defaultPage uint64 = 1
	// How many items to fetch per page
	pageLimit uint64 = 10
)

// Leaders returns a cached, alphabetized, list of leaders for the guildID. since we store these in memory, don't rely on them for non static data
func (b *Bot) Leaders(guildID uint64) ([]generated.Leader, error) {
	if leaders := b.leadersCache[guildID]; leaders != nil {
		return leaders, nil
	} else {
		leaders, err := b.Ci6ndex.GetLeaders(guildID)
		if err != nil {
			return nil, err
		}
		b.leadersCache[guildID] = leaders
		return leaders, nil
	}
}

// getNextLeader returns the alphabetically "next" leader
func (b *Bot) getNextLeader(guildID uint64, leader *generated.Leader) (*generated.Leader, error) {
	leaders, err := b.Leaders(guildID)
	if err != nil {
		return nil, err
	}

	currentIndex := -1
	for i, l := range leaders {
		if l.ID == leader.ID {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		return nil, fmt.Errorf("leader %s not found", leader.LeaderName)
	}

	// Get the next leader (cycling back to the beginning if needed)
	nextIndex := (currentIndex + 1) % len(leaders)
	return &leaders[nextIndex], nil
}

// getPrevLeader returns the alphabetically "previous" leader
func (b *Bot) getPrevLeader(guildID uint64, leader *generated.Leader) (*generated.Leader, error) {
	leaders, err := b.Leaders(guildID)
	if err != nil {
		return nil, err
	}

	currentIndex := -1
	for i, l := range leaders {
		if l.ID == leader.ID {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		return nil, fmt.Errorf("leader %s not found", leader.LeaderName)
	}

	// Get the prev leader (cycling around if needed)
	prevIndex := (currentIndex + len(leaders) - 1) % len(leaders)
	return &leaders[prevIndex], nil
}

func (b *Bot) handleManageLeadersSlashCommand() handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		guildID := e.GuildID()
		if guildID == nil {
			return fmt.Errorf("event missing guild ID")
		}
		var currentPage uint64 = 0
		components, err := b.handleManageLeaders("/leaders", guildID.String(), currentPage)

		if err != nil {
			return err
		}

		flags := discord.MessageFlagIsComponentsV2
		flags = flags.Add(discord.MessageFlagEphemeral)

		if err := e.CreateMessage(discord.MessageCreate{
			Flags:      flags,
			Components: *components,
		}); err != nil {
			slog.Error("failed to create leaders screen", "error", slog.Any("err", err))
			desc, ok := errorDescription(err)
			if ok {
				slog.Error(desc)
			}
			return err
		}
		return nil
	}
}

func (b *Bot) handleManageLeadersButtonCommand() handler.ButtonComponentHandler {
	return func(bid discord.ButtonInteractionData, e *handler.ComponentEvent) error {
		guildID := e.GuildID()
		if guildID == nil {
			return fmt.Errorf("event missing guild ID ")
		}
		var currentPage uint64
		if pageStr := e.Vars["page"]; pageStr != "" {
			if parsedPage, errConv := strconv.ParseUint(pageStr, 10, 64); errConv == nil {
				currentPage = parsedPage
			} else {
				slog.Debug("failed to parse page", "pageStr", pageStr, "error", errConv)
				currentPage = defaultPage
			}
		} else {
			slog.Debug("unable to parse page", "pageStr", pageStr)
			currentPage = defaultPage
		}

		components, err := b.handleManageLeaders(bid.CustomID(), guildID.String(), currentPage)

		if err != nil {
			return err
		}

		if err := e.UpdateMessage(discord.MessageUpdate{
			Components: components,
		}); err != nil {
			slog.Error("Failed to create leaders screen", "error", slog.Any("err", err))
			desc, ok := errorDescription(err)
			if ok {
				slog.Error(desc)
			}
			return err
		}
		return nil
	}
}

func (b *Bot) handleManageLeaders(id, guildID string, page uint64) (*[]discord.LayoutComponent, error) {
	slog.Info("handleManageLeaders", "path", id)
	gid, err := strconv.ParseUint(guildID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse guildID")
	}
	offset := (page - 1) * pageLimit
	r, err := b.leadersScreen(gid, offset, pageLimit)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (b *Bot) handleGetLeaderSlashCommand() handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		searchTerm := data.String("search")

		if searchTerm == "" {
			return errors.New("search term not provided")
		}
		return nil
		// todo:
	}
}
func (b *Bot) handleLeaderDetailsButtonCommand() handler.ButtonComponentHandler {
	return func(bid discord.ButtonInteractionData, e *handler.ComponentEvent) error {
		lid, err := strconv.ParseInt(e.Vars["leaderId"], 10, 64)
		if err != nil {
			return errors.Join(err, errors.New("failed to parse leaderId from event"))
		}
		slog.Info("handleLeaderDetails", slog.Uint64("lid", uint64(lid)), slog.String("path", bid.CustomID()))

		if e.GuildID() == nil {
			slog.Error("missing guild ID ", slog.Any("err", e))
			return errors.New("missing guild id on event")
		}
		guildId, err := parseGuildId(e.GuildID().String())
		if err != nil {
			return errors.Join(err, errors.New("failed to parse guild ID"))
		}

		// Get leader details from database
		leader, err := b.Ci6ndex.GetLeaderById(guildId, uint64(lid))
		if err != nil {
			return errors.Join(err, fmt.Errorf("failed to fetch leader with ID %d", lid))
		}

		// Create components to display leader details
		components, err := b.leaderDetailsScreen(leader, guildId)
		if err != nil {
			return err
		}

		if err := e.UpdateMessage(discord.MessageUpdate{
			Components: &components,
		}); err != nil {
			slog.Error("Failed to create leader details screen", "error", err)
			desc, ok := errorDescription(err)
			if ok {
				slog.Error(desc)
			}
		}
		return nil
	}
}

func (b *Bot) leaderDetailsScreen(leader generated.Leader, guildId uint64) ([]discord.LayoutComponent, error) {
	me, _ := b.Client.Caches.SelfUser()

	var detailsBuffer bytes.Buffer
	err := renderLeaderDetails(&detailsBuffer, leader)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to render leader details"))
	}

	prev, err := b.getPrevLeader(guildId, &leader)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to get previous leader"))
	}
	prevButton := discord.NewSecondaryButton("Previous", fmt.Sprintf("/leaders/%d", prev.ID)).
		WithEmoji(discord.ComponentEmoji{Name: "⬅️"})

	next, err := b.getNextLeader(guildId, &leader)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to get next leader"))
	}
	nextButton := discord.NewSecondaryButton("Next", fmt.Sprintf("/leaders/%d", next.ID)).
		WithEmoji(discord.ComponentEmoji{Name: "➡️"})
	layout := []discord.LayoutComponent{
		discord.NewContainer().AddComponents(
			discord.NewSection(
				discord.NewTextDisplay(detailsBuffer.String()),
			).WithAccessory(discord.NewThumbnail(me.EffectiveAvatarURL())),
			discord.NewLargeSeparator(),
			discord.NewActionRow().WithComponents(
				discord.NewPrimaryButton("Back", "/leaders"),
				prevButton,
				nextButton,
			)).
			WithAccentColor(colorSuccess),
	}

	return layout, nil
}

func renderLeaderDetails(buffer io.Writer, leader generated.Leader) error {
	md := md.NewMarkdown(buffer)

	emoji := leader.DiscordEmojiString.String

	// Create a more detailed leader profile
	mdBuilder := md.H1(emoji+" "+leader.LeaderName+" OF "+leader.CivName).
		H2f("**Tier**: %.1f", leader.Tier)
	if leader.Banned {
		mdBuilder.PlainText("\n**Status**: " + getBannedStatus(leader.Banned))
	}

	err := mdBuilder.Build()
	if err != nil {
		return errors.Join(err, errors.New("failed to build leader details markdown"))
	}

	return nil
}

func getBannedStatus(banned bool) string {
	if banned {
		return "🧑‍⚖️ Banned from drafts"
	}
	return "✅ Available for draft"
}

func (b *Bot) leadersScreen(guildId uint64, currentOffset uint64, limit uint64) ([]discord.LayoutComponent, error) {
	me, _ := b.Client.Caches.SelfUser()

	var leadersHeader bytes.Buffer
	err := renderLeadersMainScreen(&leadersHeader)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to render leaders main screen"))
	}

	// Assuming GetLeadersInRange expects a 1-based offset and a limit (number of items)
	leads, err := b.Ci6ndex.GetLeadersInRange(guildId, currentOffset, limit)
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("failed to fetch leaders with offset %d, limit %d", currentOffset, limit))
	}

	leaderRows := make([]discord.ContainerSubComponent, len(leads))
	for i, leader := range leads {
		var leaderRow bytes.Buffer
		err := md.NewMarkdown(&leaderRow).
			H1(leader.DiscordEmojiString.String + fmt.Sprintf(" %s ", leader.LeaderName)).
			Build()
		if err != nil {
			return nil, errors.Join(err, fmt.Errorf("failed to build leader: %v", leader))
		}

		leaderRoute := fmt.Sprintf("/leaders/%d", leader.ID)

		leaderRows[i] = discord.NewSection(
			discord.NewTextDisplay(
				leaderRow.String(),
			),
		).WithAccessory(discord.NewSecondaryButton(
			"Details",
			leaderRoute,
		))
	}

	layout := []discord.LayoutComponent{
		discord.NewContainer().AddComponents(
			discord.NewSection(
				discord.NewTextDisplay(leadersHeader.String()),
			).WithAccessory(discord.NewThumbnail(me.EffectiveAvatarURL())),
			discord.NewLargeSeparator()).
			AddComponents(leaderRows...).
			AddComponents(
				discord.NewLargeSeparator(),
				b.createLeaderPaginationActionRow(currentOffset, limit, len(leads)),
			).
			WithAccentColor(colorSuccess),
	}

	return layout, nil
}

func renderLeadersMainScreen(header io.Writer) error {
	err := renderLeadersHeader(header)
	if err != nil {
		return errors.Join(err, errors.New("failed to render leaders header"))
	}
	return nil
}

func renderLeadersHeader(header io.Writer) error {
	return md.NewMarkdown(header).H1("Browse Civs").
		PlainText("View statistics about each leader, and optionally update your personal rank of them.").
		Build()
}

func (b *Bot) createLeaderPaginationActionRow(currentOffset, limit uint64, numLeadsFetched int) discord.ActionRowComponent {
	prevPage := (currentOffset / limit)
	prevButton := discord.NewSecondaryButton("Previous", fmt.Sprintf("/leaders/page/%d", prevPage)).
		WithEmoji(discord.ComponentEmoji{Name: "⬅️"})
	if currentOffset <= 1 {
		prevButton = prevButton.WithDisabled(true)
	}

	nextPage := (currentOffset / limit) + 2

	nextButton := discord.NewSecondaryButton("Next", fmt.Sprintf("/leaders/page/%d", nextPage)).
		WithEmoji(discord.ComponentEmoji{Name: "➡️"})
	if numLeadsFetched < int(limit) {
		nextButton = nextButton.WithDisabled(true)
	}
	// Back Button (to draft menu)
	backToDraftButton := discord.NewPrimaryButton("Drafts", "/draft").WithEmoji(discord.ComponentEmoji{
		Name: crossedSwords,
	})

	return discord.NewActionRow().WithComponents(backToDraftButton, prevButton, nextButton)
}
