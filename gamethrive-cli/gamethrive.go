package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"../gamethrive"
)

var (
	PlayerFlagSet           *flag.FlagSet
	PlayerJsonPathFlag      *string
	PlayerAppIdFlag         *string
	PlayerIdFlag            *string
	PlayerDeviceTypeFlag    *string
	PlayerIdentifierFlag    *string
	PlayerLanguageFlag      *string
	PlayerTimezoneFlag      *int
	PlayerDeviceModelFlag   *string
	PlayerDeviceOSFlag      *string
	PlayerGameVerionFlag    *string
	PlayerAdvertisingIdFlag *string
	PlayerSessionCountFlag  *int
	PlayerTagsFlag          *string
	PlayerAmountSpentFlag   *float64
	PlayerCreatedAtFlag     *int
	PlayerLastActiveFlag    *int
	PlayerPlaytimeFlag      *int

	PlayerAmountFlagSet    *flag.FlagSet
	PlayerAmountIdFlag     *string
	PlayerAmountAmountFlag *float64

	PlayerPlaytimeFlagSet   *flag.FlagSet
	PlayerPlaytimeIdFlag    *string
	PlayerPlaytimeStateFlag *string
	PlayerPlaytimeTimeFlag  *int

	NotificationFlagSet                   *flag.FlagSet
	NotificationJsonPathFlag              *string
	NotificationAuthPathFlag              *string
	NotificationAppIdFlag                 *string
	NotificationIsIOSFlag                 *bool
	NotificationIsAndroidFlag             *bool
	NotificationContentsFlag              *string
	NotificationIncludedSegmentsFlag      *string
	NotificationExcludedSegmentsFlag      *string
	NotificationIncludedPlayerIdsFlag     *string
	NotificationIncludedIOSTokensFlag     *string
	NotificationIncludedAndroidRegIdsFlag *string
	NotificationIOSBadgeTypeFlag          *string
	NotificationIOSBadgeCountFlag         *int
	NotificationIOSSoundFlag              *string
	NotificationAndroidSoundFlag          *string
	NotificationDataFlag                  *string
	NotificationURLFlag                   *string
	NotificationSendAfterFlag             *string
	NotificationSendUserActiveTimeFlag    *bool

	NotificationOpenFlagSet    *flag.FlagSet
	NotificationOpenIdFlag     *string
	NotificationOpenAppIdFlag  *string
	NotificationOpenOpenedFlag *bool
)

func init() {
	PlayerFlagSet = flag.NewFlagSet("player", flag.ContinueOnError)
	PlayerJsonPathFlag = PlayerFlagSet.String("json", "", "Read player info from a json file")
	PlayerAppIdFlag = PlayerFlagSet.String("app_id", "", "Your GameThrive's application key")
	PlayerIdFlag = PlayerFlagSet.String("id", "", "Gamethrive identifier of the player")
	PlayerDeviceTypeFlag = PlayerFlagSet.String("device_type", "ios", `"ios", "android" or "amazon"`)
	PlayerIdentifierFlag = PlayerFlagSet.String("identifier", "", "Push notification identifier from Google or Apple")
	PlayerLanguageFlag = PlayerFlagSet.String("language", "", "Language code. Typically lower case two letters, except for chinese")
	PlayerTimezoneFlag = PlayerFlagSet.Int("timezone", 0, "Number of seconds away from GMT")
	PlayerDeviceModelFlag = PlayerFlagSet.String("device_model", "", "Device model")
	PlayerDeviceOSFlag = PlayerFlagSet.String("device_os", "", "Device operating system version")
	PlayerGameVerionFlag = PlayerFlagSet.String("game_version", "", "Version of the game")
	PlayerAdvertisingIdFlag = PlayerFlagSet.String("ad_id", "", "Advertising id for Android devices and identifierForVendor for iOS devices")
	PlayerSessionCountFlag = PlayerFlagSet.Int("session_count", 1, "Number of times the player has played the game, defaults to 1")
	PlayerTagsFlag = PlayerFlagSet.String("tags", "{}", "Custom tags for the player (a json string)")
	PlayerAmountSpentFlag = PlayerFlagSet.Float64("amount_spent", 0.0, "Amount the player has spent in USD, up to two decimal places")
	PlayerCreatedAtFlag = PlayerFlagSet.Int("created_at", 0, "Unixtime when the player joined the game")
	PlayerLastActiveFlag = PlayerFlagSet.Int("last_active", 0, "Unixtime when the player was last active")
	PlayerPlaytimeFlag = PlayerFlagSet.Int("playtime", 0, "Seconds player was running your app.")

	PlayerAmountFlagSet = flag.NewFlagSet("player amount", flag.ContinueOnError)
	PlayerAmountIdFlag = PlayerAmountFlagSet.String("id", "", "Gamethrive identifier of the player")
	PlayerAmountAmountFlag = PlayerAmountFlagSet.Float64("amount", 0.0, "New amount in USD, up to two decimal places")

	PlayerPlaytimeFlagSet = flag.NewFlagSet("player playtime", flag.ContinueOnError)
	PlayerPlaytimeIdFlag = PlayerPlaytimeFlagSet.String("id", "", "Gamethrive identifier of the player")
	PlayerPlaytimeStateFlag = PlayerPlaytimeFlagSet.String("state", "", "Required to indicate we are incrementing")
	PlayerPlaytimeTimeFlag = PlayerPlaytimeFlagSet.Int("active_time", 0, "Number of seconds player was running your app")

	NotificationFlagSet = flag.NewFlagSet("notification", flag.ContinueOnError)
	NotificationJsonPathFlag = NotificationFlagSet.String("json", "", "Read notification info from a json file")
	NotificationAuthPathFlag = NotificationFlagSet.String("auth", "", `Your "API Auth Key" on the GameThrive Application Settings page`)
	NotificationAppIdFlag = NotificationFlagSet.String("app_id", "", "Your GameThrive's application key")
	NotificationIsIOSFlag = NotificationFlagSet.Bool("ios", false, "Send notification to iOS players")
	NotificationIsAndroidFlag = NotificationFlagSet.Bool("android", false, "Send notification to Android players")
	NotificationContentsFlag = NotificationFlagSet.String("contents", `{"en":""}`, "Message contents to send to players, \"en\" (English) is required")
	NotificationIncludedSegmentsFlag = NotificationFlagSet.String("included_segments", "", "Names of segments to send the message to (separated by commas)")
	NotificationExcludedSegmentsFlag = NotificationFlagSet.String("excluded_segments", "", "Names of segments to exclude players from (separated by commas)")
	NotificationIncludedPlayerIdsFlag = NotificationFlagSet.String("include_player_ids", "", "Specific players to send your notification to (separated by commas)")
	NotificationIncludedIOSTokensFlag = NotificationFlagSet.String("include_ios_tokens", "", "Specific iOS device tokens to send the notification to (separated by commas)")
	NotificationIncludedAndroidRegIdsFlag = NotificationFlagSet.String("include_android_reg_ids", "", "Specific Android registration ids to send the notification to (separated by commas)")
	NotificationIOSBadgeTypeFlag = NotificationFlagSet.String("ios_badgeType", "none", `Options are: "none", "setto", or "increase"`)
	NotificationIOSBadgeCountFlag = NotificationFlagSet.Int("ios_badgeCount", 0, "Sets or increases the badge icon on the device")
	NotificationIOSSoundFlag = NotificationFlagSet.String("ios_sound", "", "Sound file that is included in your app to play")
	NotificationAndroidSoundFlag = NotificationFlagSet.String("android_sound", "", "Sound file that is included in your app to play")
	NotificationDataFlag = NotificationFlagSet.String("data", "", "Custom key value pair hash that you can programmatically read in your app's code (as json string)")
	NotificationURLFlag = NotificationFlagSet.String("url", "", "When the player opens the notification their web browser will open this url")
	NotificationSendAfterFlag = NotificationFlagSet.String("send_after", "", `Schedule notification for future delivery (Format: "Mon Jan 02 2006 15:04:05 MST-0700")`)
	NotificationSendUserActiveTimeFlag = NotificationFlagSet.Bool("send_at_user_active_time", false, "Sends your notification at the time of day the user last opened your app")

	NotificationOpenFlagSet = flag.NewFlagSet("notification open", flag.ContinueOnError)
	NotificationOpenIdFlag = NotificationOpenFlagSet.String("id", "", "Identifier of the notification")
	NotificationOpenAppIdFlag = NotificationOpenFlagSet.String("app_id", "", "Your GameThrive's application key")
	NotificationOpenOpenedFlag = NotificationOpenFlagSet.Bool("opened", true, "Required to indicate the notification was openned")
}

func main() {
	action := NewAction(os.Args[1:])
	mux := map[string]interface{}{
		"players": map[string]interface{}{
			"new": map[string]interface{}{
				"handler": Handler(PlayersNew),
				"usage":   "Creates a player",
			},
			"update": map[string]interface{}{
				"handler": Handler(PlayersUpdate),
				"usage":   "Updates player attributes",
			},
			"amount": map[string]interface{}{
				"handler": Handler(PlayerUpdateAmount),
				"usage":   "Updates player's amount",
			},
			"session": map[string]interface{}{
				"handler": Handler(PlayerSession),
				"usage":   "Updates and increments player session count",
			},
			"playtime": map[string]interface{}{
				"handler": Handler(PlayerPlaytime),
				"usage":   "Increment the player's total playtime",
			},
		},
		"notifications": map[string]interface{}{
			"new": map[string]interface{}{
				"handler": Handler(NotificationNew),
				"usage":   "Create and deliver a new a Notification",
			},
			"open": map[string]interface{}{
				"handler": Handler(NotificationOpen),
				"usage":   "Track that a push notification was opened",
			},
		},
		"help": map[string]interface{}{
			"players": map[string]interface{}{
				"new": map[string]interface{}{
					"handler": Handler(HelpPlayersNew),
				},
				"update": map[string]interface{}{
					"handler": Handler(HelpPlayersUpdate),
				},
				"amount": map[string]interface{}{
					"handler": Handler(HelpPlayersUpdateAmount),
				},
				"session": map[string]interface{}{
					"handler": Handler(HelpPlayerSession),
				},
				"playtime": map[string]interface{}{
					"handler": Handler(HelpPlayerPlaytime),
				},
			},
			"notifications": map[string]interface{}{
				"new": map[string]interface{}{
					"handler": Handler(HelpNotificationNew),
				},
				"open": map[string]interface{}{
					"handler": Handler(HelpNotificationOpen),
				},
			},
			"usage": "Shows usage about each command. Example: help players new",
		},
	}
	handler := GetHandler(mux, action)
	if handler == nil {
		fmt.Printf("ERROR\n%s\n", GetUsage(mux))
		return
	}
	handler(os.Args[len(action)+1:]...)
}

type Handler func(args ...string)

func NewAction(args []string) []string {
	acs := []string{}
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			break
		}
		acs = append(acs, a)
	}
	return acs
}

func GetHandler(mux map[string]interface{}, action []string) Handler {
	if len(action) <= 0 {
		return nil
	}
	leaf, ok := mux[action[0]]
	if !ok {
		return nil
	}
	lmux, ok := leaf.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(action) > 1 {
		return GetHandler(lmux, action[1:])
	}
	handler, ok := lmux["handler"].(Handler)
	if !ok {
		return nil
	}
	return handler
}

func GetUsage(mux map[string]interface{}) string {
	return joinAlign(visitUsage(mux))
}

func PlayersNew(args ...string) {
	PlayerFlagSet.Parse(args)
	player, err := currentPlayer()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	c := gamethrive.NewClient(nil)
	err = c.Players.New(player)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	fmt.Printf("Player created correctly. Player id is: \"%s\"\n", player.Id)
}

func HelpPlayersNew(args ...string) {
	fmt.Println("Creates a new player.")
	PlayerFlagSet.PrintDefaults()
}

func PlayersUpdate(args ...string) {
	PlayerFlagSet.Parse(args)
	player, err := currentPlayer()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	c := gamethrive.NewClient(nil)
	err = c.Players.Update(player)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
}

func HelpPlayersUpdate(args ...string) {
	fmt.Println("Updates player attributes.")
	fmt.Println("Note: Updating tags will append to the player's existing tags.")
	fmt.Println("      To remove an existing tag, update the tag with a value of a blank string.")
	PlayerFlagSet.PrintDefaults()
}

func PlayerUpdateAmount(args ...string) {
	PlayerAmountFlagSet.Parse(args)
	if len(*PlayerAmountIdFlag) <= 0 {
		fmt.Println("Error: id flag is requried")
		return
	}
	c := gamethrive.NewClient(nil)
	err := c.Players.UpdateAmount(*PlayerAmountIdFlag, *PlayerAmountAmountFlag)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
}

func HelpPlayersUpdateAmount(args ...string) {
	fmt.Println("Increment player's total amount spent.")
	PlayerAmountFlagSet.PrintDefaults()
}

func PlayerSession(args ...string) {
	PlayerFlagSet.Parse(args)
	player, err := currentPlayer()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	c := gamethrive.NewClient(nil)
	err = c.Players.Session(player)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
}

func HelpPlayerSession(args ...string) {
	fmt.Println("Updates any player details that may have changed as well as incrementing the player's session count.")
	PlayerFlagSet.PrintDefaults()
}

func PlayerPlaytime(args ...string) {
	PlayerPlaytimeFlagSet.Parse(args)
	if len(*PlayerPlaytimeIdFlag) <= 0 {
		fmt.Println("Error: id flag is requried")
		return
	}
	c := gamethrive.NewClient(nil)
	state := stringToPlayState(*PlayerPlaytimeStateFlag)
	err := c.Players.Playtime(*PlayerPlaytimeIdFlag, state, *PlayerPlaytimeTimeFlag)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
}

func HelpPlayerPlaytime(args ...string) {
	fmt.Println("Increment the player's total playtime")
	PlayerPlaytimeFlagSet.PrintDefaults()
}

func NotificationNew(args ...string) {
	NotificationFlagSet.Parse(args)
	notification, err := currentNotification()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	c := gamethrive.NewClient(nil)
	d, err := c.Notifications.New(notification, *NotificationAuthPathFlag)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	fmt.Printf("Notification (%s) created sucessfully. Target: %d players\n", notification.Id, d)
}

func HelpNotificationNew(args ...string) {
	fmt.Println("Create and deliver a new a Notification")
	NotificationFlagSet.PrintDefaults()
}

func NotificationOpen(args ...string) {
	NotificationOpenFlagSet.Parse(args)
	if len(*NotificationOpenIdFlag) <= 0 {
		fmt.Println("Error: id flag is requried")
		return
	}
	if len(*NotificationOpenAppIdFlag) <= 0 {
		fmt.Println("Error: app_id flag is requried")
		return
	}
	c := gamethrive.NewClient(nil)
	notification := gamethrive.Notification{
		Id:    *NotificationOpenIdFlag,
		AppId: *NotificationOpenAppIdFlag,
	}
	err := c.Notifications.Open(&notification, *NotificationOpenOpenedFlag)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}

func HelpNotificationOpen(args ...string) {
	fmt.Println("Track that a push notification was opened")
	NotificationOpenFlagSet.PrintDefaults()
}

func currentPlayer() (*gamethrive.Player, error) {
	player := new(gamethrive.Player)
	if len(*PlayerJsonPathFlag) > 0 {
		file, err := os.Open(*PlayerJsonPathFlag)
		defer file.Close()
		if err != nil {
			return nil, err
		}
		err = json.NewDecoder(file).Decode(player)
		return player, err
	}
	player.AppId = *PlayerAppIdFlag
	player.Id = *PlayerIdFlag
	player.DeviceType = stringToDeviceType(*PlayerDeviceTypeFlag)
	player.Identifier = *PlayerIdentifierFlag
	player.Language = *PlayerLanguageFlag
	player.Timezone = *PlayerTimezoneFlag
	player.DeviceModel = *PlayerDeviceModelFlag
	player.DeviceOS = *PlayerDeviceOSFlag
	player.GameVersion = *PlayerGameVerionFlag
	player.AdvertisingId = *PlayerAdvertisingIdFlag
	player.SessionCount = *PlayerSessionCountFlag
	tags, err := currentPlayerTags()
	if err != nil {
		return nil, err
	}
	player.Tags = tags
	player.AmountSpent = *PlayerAmountSpentFlag
	player.CreatedAt = *PlayerCreatedAtFlag
	player.LastActive = *PlayerLastActiveFlag
	player.Playtime = *PlayerPlaytimeFlag
	return player, nil
}

func currentPlayerTags() (m map[string]string, err error) {
	buffer := ioutil.NopCloser(strings.NewReader(*PlayerTagsFlag))
	err = json.NewDecoder(buffer).Decode(&m)
	return
}

func currentNotification() (*gamethrive.Notification, error) {
	notification := new(gamethrive.Notification)
	if len(*NotificationJsonPathFlag) > 0 {
		file, err := os.Open(*NotificationJsonPathFlag)
		defer file.Close()
		if err != nil {
			return nil, err
		}
		err = json.NewDecoder(file).Decode(&notification)
		return notification, err
	}
	notification.AppId = *NotificationAppIdFlag
	notification.IsIOS = *NotificationIsIOSFlag
	notification.IsAndroid = *NotificationIsAndroidFlag
	notification.Contents = currentNotificationContents()
	if len(*NotificationIncludedSegmentsFlag) > 0 {
		notification.IncludedSegments = strings.Split(*NotificationIncludedSegmentsFlag, ",")
	}
	if len(*NotificationExcludedSegmentsFlag) > 0 {
		notification.ExcludedSegments = strings.Split(*NotificationExcludedSegmentsFlag, ",")
	}
	if len(*NotificationIncludedPlayerIdsFlag) > 0 {
		notification.IncludedPlayerIds = strings.Split(*NotificationIncludedPlayerIdsFlag, ",")
	}
	if len(*NotificationIncludedIOSTokensFlag) > 0 {
		notification.IncludedIOSTokens = strings.Split(*NotificationIncludedIOSTokensFlag, ",")
	}
	if len(*NotificationIncludedAndroidRegIdsFlag) > 0 {
		notification.IncludedAndroidRegIds = strings.Split(*NotificationIncludedAndroidRegIdsFlag, ",")
	}
	notification.IOSBadgeType = stringToBadgeType(*NotificationIOSBadgeTypeFlag)
	notification.IOSBadgeCount = *NotificationIOSBadgeCountFlag
	notification.IOSSound = *NotificationIOSSoundFlag
	notification.AndroidSound = *NotificationAndroidSoundFlag
	notification.Data = currentNotificationData()
	notification.URL = *NotificationURLFlag
	if len(*NotificationSendAfterFlag) > 0 {
		t, err := time.Parse("Mon Jan 02 2006 15:04:05 MST-0700", *NotificationSendAfterFlag)
		if err != nil {
			return nil, err
		}
		notification.SendAfter = &t
	}
	notification.SendUserActiveTime = *NotificationSendUserActiveTimeFlag
	return notification, nil
}

func currentNotificationContents() (c map[string]string) {
	buffer := ioutil.NopCloser(strings.NewReader(*NotificationContentsFlag))
	json.NewDecoder(buffer).Decode(&c)
	return
}

func currentNotificationData() (c map[string]string) {
	buffer := ioutil.NopCloser(strings.NewReader(*NotificationDataFlag))
	json.NewDecoder(buffer).Decode(&c)
	return
}

func stringToDeviceType(str string) gamethrive.DeviceType {
	switch str {
	case "android":
		return gamethrive.Android
	case "amazon":
		return gamethrive.Amazon
	default:
		return gamethrive.IOS
	}
}

func stringToPlayState(str string) gamethrive.PlaytimeState {
	switch str {
	case "suspend":
		return gamethrive.Suspend
	case "resume":
		return gamethrive.Resume
	default:
		return gamethrive.Ping
	}
}

func stringToBadgeType(str string) gamethrive.BadgeType {
	switch str {
	case "setto":
		return gamethrive.SetTo
	case "increase":
		return gamethrive.Increase
	default:
		return gamethrive.None
	}
}

func visitUsage(mux map[string]interface{}) []string {
	if usage, ok := mux["usage"].(string); ok {
		return []string{"\r" + usage}
	}
	var lines []string
	for k, v := range mux {
		leaf, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		sublines := visitUsage(leaf)
		if len(sublines) <= 0 {
			continue
		}
		lines = append(lines, appendLines(k+" ", sublines)...)
	}
	return lines
}

func appendLines(prefix string, lines []string) []string {
	newlines := []string{}
	for _, line := range lines {
		newlines = append(newlines, prefix+line)
	}
	return newlines
}

func joinAlign(lines []string) string {
	column := returnColumn(lines)
	newlines := []string{}
	for _, l := range lines {
		newlines = append(newlines, adjustColumn(l, column))
	}
	return strings.Join(newlines, "\n")
}

func adjustColumn(line string, column int) string {
	gs := strings.Split(line, "\r")
	if len(gs) != 2 {
		return line
	}
	return gs[0] + strings.Repeat(" ", column-len(gs[0])) + "\t" + gs[1]
}

func returnColumn(lines []string) (c int) {
	for _, l := range lines {
		gs := strings.Split(l, "\r")
		if len(gs) <= 1 {
			continue
		}
		cc := lens(gs[:(len(gs) - 1)]...)
		if cc > c {
			c = cc
		}
	}
	return
}

func lens(strs ...string) (l int) {
	for _, s := range strs {
		l += len(s)
	}
	return
}
