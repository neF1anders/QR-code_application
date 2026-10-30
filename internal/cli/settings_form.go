package cli

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"time"

	"github.com/charmbracelet/huh"

	"github.com/neF1anders/QR-code_application/internal/domain"
)

// customTTL marks the "enter seconds manually" choice in the TTL list.
const customTTL time.Duration = -1

var ttlPresets = []time.Duration{0, time.Minute, 5 * time.Minute, time.Hour}

var sizePresets = []int{256, 384, 512, 768}

// EditSettingsInteractively shows dropdown lists in the terminal.
func EditSettingsInteractively(current domain.Settings) (domain.Settings, error) {
	ttlChoice := ttlChoiceFor(current.TTL)
	customSeconds := strconv.FormatInt(int64(current.TTL/time.Second), 10)
	size := current.SizePx
	recovery := current.Recovery

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[time.Duration]().
				Title("Image lifetime (TTL)").
				Description("How long the QR image is kept on disk after you close the window").
				Options(ttlOptions()...).
				Value(&ttlChoice),
			huh.NewSelect[int]().
				Title("Image size").
				Options(sizeOptions(current.SizePx)...).
				Value(&size),
			huh.NewSelect[domain.RecoveryLevel]().
				Title("Error correction").
				Description("Higher survives more damage, lower fits more text").
				Options(recoveryOptions()...).
				Value(&recovery),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("TTL in seconds").
				Value(&customSeconds).
				Validate(func(input string) error {
					_, err := parseSeconds(input)
					return err
				}),
		).WithHideFunc(func() bool { return ttlChoice != customTTL }),
	)
	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return domain.Settings{}, ErrEditCanceled
		}
		return domain.Settings{}, err
	}

	ttl := ttlChoice
	if ttlChoice == customTTL {
		var err error
		if ttl, err = parseSeconds(customSeconds); err != nil {
			return domain.Settings{}, err
		}
	}
	return domain.Settings{TTL: ttl, SizePx: size, Recovery: recovery}, nil
}

func ttlChoiceFor(ttl time.Duration) time.Duration {
	for _, preset := range ttlPresets {
		if ttl == preset {
			return preset
		}
	}
	return customTTL
}

func ttlOptions() []huh.Option[time.Duration] {
	options := []huh.Option[time.Duration]{huh.NewOption("Delete when the window closes", time.Duration(0))}
	for _, preset := range ttlPresets[1:] {
		options = append(options, huh.NewOption("Keep for "+preset.String(), preset))
	}
	return append(options, huh.NewOption("Custom (seconds)…", customTTL))
}

// sizeOptions includes the current size so a hand-edited config value is not lost.
func sizeOptions(current int) []huh.Option[int] {
	sizes := sizePresets
	if !slices.Contains(sizes, current) {
		sizes = append([]int{current}, sizes...)
	}
	options := make([]huh.Option[int], 0, len(sizes))
	for _, size := range sizes {
		options = append(options, huh.NewOption(fmt.Sprintf("%d px", size), size))
	}
	return options
}

func recoveryOptions() []huh.Option[domain.RecoveryLevel] {
	options := make([]huh.Option[domain.RecoveryLevel], 0, len(domain.RecoveryLevels))
	for _, level := range domain.RecoveryLevels {
		options = append(options, huh.NewOption(string(level), level))
	}
	return options
}

func parseSeconds(input string) (time.Duration, error) {
	seconds, err := strconv.ParseInt(input, 10, 64)
	if err != nil || seconds < 0 || seconds > math.MaxInt64/int64(time.Second) {
		return 0, fmt.Errorf("enter a whole number of seconds, 0 or more")
	}
	return time.Duration(seconds) * time.Second, nil
}
