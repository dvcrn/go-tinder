package tinder

import "time"

type Meta struct {
	Status int `json:"status"`
}

type Badge struct {
	Type string `json:"type"`
}

type Person struct {
	ID        string   `json:"_id"`
	Badges    []*Badge `json:"badges"`
	Bio       string   `json:"bio"`
	BirthDate string   `json:"birth_date"`
	Gender    int      `json:"gender"`
	Name      string   `json:"name"`
	PingTime  string   `json:"ping_time"`
	Photos    []*Photo `json:"photos"`
}

type Message struct {
	ID          string    `json:"_id"`
	MatchID     string    `json:"match_id"`
	SentDate    time.Time `json:"sent_date"`
	CreatedDate time.Time `json:"created_date"`
	Message     string    `json:"message"`
	To          string    `json:"to"`
	From        string    `json:"from"`
	Timestamp   int       `json:"timestamp"`
	MatchId     string    `json:"matchId"`
}

type Match struct {
	Seen                    *Seen      `json:"seen"`
	ID                      string     `json:"_id"`
	Closed                  bool       `json:"closed"`
	CommonFriendCount       int        `json:"common_friend_count"`
	CommonLikeCount         int        `json:"common_like_count"`
	CreatedDate             time.Time  `json:"created_date"`
	Dead                    bool       `json:"dead"`
	LastActivityDate        time.Time  `json:"last_activity_date"`
	MessageCount            int        `json:"message_count"`
	Messages                []*Message `json:"messages"`
	Participants            []string   `json:"participants"`
	Pending                 bool       `json:"pending"`
	IsSuperLike             bool       `json:"is_super_like"`
	IsBoostMatch            bool       `json:"is_boost_match"`
	IsSuperBoostMatch       bool       `json:"is_super_boost_match"`
	IsPrimetimeBoostMatch   bool       `json:"is_primetime_boost_match"`
	IsExperiencesMatch      bool       `json:"is_experiences_match"`
	IsFastMatch             bool       `json:"is_fast_match"`
	IsPreferencesMatch      bool       `json:"is_preferences_match"`
	IsMatchmakerMatch       bool       `json:"is_matchmaker_match"`
	IsLetsMeetMatch         bool       `json:"is_lets_meet_match"`
	IsOpener                bool       `json:"is_opener"`
	HasShownInitialInterest bool       `json:"has_shown_initial_interest"`
	Person                  *Person    `json:"person"`
	Following               bool       `json:"following"`
	FollowingMoments        bool       `json:"following_moments"`
	ReadReceipt             struct {
		Enabled bool `json:"enabled"`
	} `json:"readreceipt"`
	LikedContent struct {
		ByCloser struct {
			UserID      string `json:"user_id"`
			Type        string `json:"type"`
			IsSwipeNote bool   `json:"is_swipe_note"`
		} `json:"by_closer"`
	} `json:"liked_content"`
	SubscriptionTier string `json:"subscription_tier"`
}

type UpdateResponse struct {
	Matches           []*Match  `json:"matches"`
	Blocks            []string  `json:"blocks"`
	Inbox             []any     `json:"inbox"`
	LikedMessages     []any     `json:"liked_messages"`
	HarassingMessages []any     `json:"harassing_messages"`
	Lists             []any     `json:"lists"`
	Goingout          []any     `json:"goingout"`
	DeletedLists      []any     `json:"deleted_lists"`
	Squads            []any     `json:"squads"`
	LastActivityDate  time.Time `json:"last_activity_date"`
	PollInterval      struct {
		Standard   int `json:"standard"`
		Persistent int `json:"persistent"`
	} `json:"poll_interval"`
	SelfieVerification struct {
		Status string `json:"status"`
	} `json:"selfie_verification"`
}

type Seen struct {
	MatchSeen     bool   `json:"match_seen"`
	LastSeenMsgID string `json:"last_seen_msg_id"`
}

type Data struct {
	Matches  *[]*Match    `json:"matches"`
	Messages *[]*Message  `json:"messages"`
	User     *UserProfile `json:"user"`
}

type Response struct {
	Meta *Meta `json:"meta"`
	Data *Data `json:"data"`
}

type UserProfileResponse struct {
	Status  int          `json:"status"`
	Results *UserProfile `json:"results"`
}

type UserProfile struct {
	ID                       string                  `json:"_id"`
	SNumber                  int64                   `json:"s_number"`
	IsTinderU                bool                    `json:"is_tinder_u"`
	CommonFriends            []any                   `json:"common_friends"`
	CommonFriendCount        int                     `json:"common_friend_count"`
	SpotifyTopArtists        []any                   `json:"spotify_top_artists"`
	DistanceMi               int                     `json:"distance_mi"`
	ConnectionCount          int                     `json:"connection_count"`
	CommonConnections        []any                   `json:"common_connections"`
	Bio                      string                  `json:"bio"`
	BirthDate                time.Time               `json:"birth_date"`
	Name                     string                  `json:"name"`
	Jobs                     []any                   `json:"jobs"`
	Schools                  []any                   `json:"schools"`
	Teasers                  []any                   `json:"teasers"`
	Gender                   int                     `json:"gender"`
	ShowGenderOnProfile      bool                    `json:"show_gender_on_profile"`
	SexualOrientations       []*SexualOrientation    `json:"sexual_orientations"`
	BirthDateInfo            string                  `json:"birth_date_info"`
	PingTime                 string                  `json:"ping_time"`
	ShowOrientationOnProfile bool                    `json:"show_orientation_on_profile"`
	Badges                   []*Badge                `json:"badges"`
	Photos                   []*Photo                `json:"photos"`
	CommonLikes              []any                   `json:"common_likes"`
	CommonLikeCount          int                     `json:"common_like_count"`
	CommonInterests          []any                   `json:"common_interests"`
	SelectedDescriptors      []*Descriptor           `json:"selected_descriptors"`
	RelationshipIntent       *RelationshipIntent     `json:"relationship_intent"`
	SparksQuizzes            []*SparksQuiz           `json:"sparks_quizzes"`
	UserPrompts              *UserPrompts            `json:"user_prompts"`
	ProfileDetailContent     []*ProfileDetailContent `json:"profile_detail_content"`
}

type SexualOrientation struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Photo struct {
	ID              string           `json:"id"`
	CropInfo        *CropInfo        `json:"crop_info"`
	URL             string           `json:"url"`
	ProcessedFiles  []*ProcessedFile `json:"processedFiles"`
	ProcessedVideos []any            `json:"processedVideos"`
	FileName        string           `json:"fileName"`
	Extension       string           `json:"extension"`
	Assets          []*Asset         `json:"assets"`
	MediaType       string           `json:"media_type"`
}

type CropInfo struct {
	User                *UserCrop `json:"user"`
	Algo                *AlgoCrop `json:"algo"`
	ProcessedByBullseye bool      `json:"processed_by_bullseye"`
	UserCustomized      bool      `json:"user_customized"`
	Faces               []*Face   `json:"faces"`
}

type UserCrop struct {
	WidthPct   float64 `json:"width_pct"`
	XOffsetPct float64 `json:"x_offset_pct"`
	HeightPct  float64 `json:"height_pct"`
	YOffsetPct float64 `json:"y_offset_pct"`
}

type AlgoCrop struct {
	WidthPct   float64 `json:"width_pct"`
	XOffsetPct float64 `json:"x_offset_pct"`
	HeightPct  float64 `json:"height_pct"`
	YOffsetPct float64 `json:"y_offset_pct"`
}

type Face struct {
	Algo                  *AlgoCrop `json:"algo"`
	BoundingBoxPercentage float64   `json:"bounding_box_percentage"`
}

type ProcessedFile struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type Asset struct {
	URL       string `json:"url"`
	AssetType string `json:"asset_type"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

type Descriptor struct {
	ID                  string               `json:"id"`
	Name                string               `json:"name"`
	Prompt              string               `json:"prompt"`
	Type                string               `json:"type"`
	IconURL             string               `json:"icon_url"`
	IconURLs            []*IconURL           `json:"icon_urls"`
	ChoiceSelections    []*ChoiceSelection   `json:"choice_selections"`
	SectionID           string               `json:"section_id"`
	SectionName         string               `json:"section_name"`
	MeasurableSelection *MeasurableSelection `json:"measurable_selection"`
}

type IconURL struct {
	URL     string `json:"url"`
	Quality string `json:"quality"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
}

type ChoiceSelection struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	MatchGroupKey string `json:"match_group_key"`
}

type MeasurableSelection struct {
	Value         int    `json:"value"`
	Min           int    `json:"min"`
	Max           int    `json:"max"`
	UnitOfMeasure string `json:"unit_of_measure"`
}

type RelationshipIntent struct {
	DescriptorChoiceID string        `json:"descriptor_choice_id"`
	Emoji              string        `json:"emoji"`
	ImageURL           string        `json:"image_url"`
	TitleText          string        `json:"title_text"`
	BodyText           string        `json:"body_text"`
	Style              string        `json:"style"`
	HiddenIntent       *HiddenIntent `json:"hidden_intent"`
	TappedAction       *TappedAction `json:"tapped_action"`
}

type HiddenIntent struct {
	Emoji     string `json:"emoji"`
	ImageURL  string `json:"image_url"`
	TitleText string `json:"title_text"`
	BodyText  string `json:"body_text"`
}

type TappedAction struct {
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	QueryParams map[string]string `json:"query_params"`
}

type SparksQuiz struct {
	Quizzes          []*Quiz          `json:"quizzes"`
	SectionID        string           `json:"section_id"`
	SectionName      string           `json:"section_name"`
	SimilarityScore  *SimilarityScore `json:"similarity_score"`
	LockedButtonText string           `json:"locked_button_text"`
}

type Quiz struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	Answers         []string         `json:"answers"`
	AnswerDetails   []*AnswerDetail  `json:"answer_details"`
	ImageURL        string           `json:"image_url"`
	LockedImageURL  string           `json:"locked_image_url"`
	SimilarityScore *SimilarityScore `json:"similarity_score"`
}

type AnswerDetail struct {
	Emoji      string `json:"emoji"`
	PromptText string `json:"prompt_text"`
	AnswerID   string `json:"answer_id"`
	AnswerText string `json:"answer_text"`
}

type SimilarityScore struct {
	Text           string `json:"text"`
	Style          string `json:"style"`
	ImageURL       string `json:"image_url"`
	TitleText      string `json:"title_text"`
	PromptText     string `json:"prompt_text"`
	YouPreferText  string `json:"you_prefer_text"`
	TheyPreferText string `json:"they_prefer_text"`
}

type UserPrompts struct {
	SectionName string     `json:"section_name"`
	Prompts     []*Prompt  `json:"prompts"`
	AddPrompt   *AddPrompt `json:"add_prompt"`
}

type Prompt struct {
	ID           string `json:"id"`
	QuestionText string `json:"question_text"`
	AnswerID     string `json:"answer_id"`
	AnswerText   string `json:"answer_text"`
	ImageURL     string `json:"image_url"`
}

type AddPrompt struct {
	Text         string        `json:"text"`
	TappedAction *TappedAction `json:"tapped_action"`
}

type ProfileDetailContent struct {
	Content       []any  `json:"content"`
	PageContentID string `json:"page_content_id"`
}
