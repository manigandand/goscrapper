package main

// KaggleUser user data csv
type KaggleUser struct {
	Id int `json:"id"`
	// UserID          int    `json:"UserName"`
	UserName        string `json:"user_name"`
	DisplayName     string `json:"display_name"`
	RegisterDate    string `json:"register_date"`
	PerformanceTier string `json:"performance_tier"`
	IsVisited       bool   `json:"is_visited" sql:"default: false"`
}

// Kaggle holds the kaggle user data
type Kaggle struct {
	UserID                  int    `json:"userId"`
	DisplayName             string `json:"displayName"`
	Country                 string `json:"country"`
	Region                  string `json:"region"`
	City                    string `json:"city"`
	GitHubUserName          string `json:"gitHubUserName"`
	TwitterUserName         string `json:"twitterUserName"`
	LinkedInURL             string `json:"linkedInUrl"`
	WebsiteURL              string `json:"websiteUrl"`
	UserURL                 string `json:"userUrl"`
	UserName                string `json:"userName"`
	Occupation              string `json:"occupation"`
	Organization            string `json:"organization"`
	Bio                     string `json:"bio"`
	UserAvatarURL           string `json:"userAvatarUrl"`
	Email                   string `json:"email"`
	UserLastActive          string `json:"userLastActive"`
	UserJoinDate            string `json:"userJoinDate"`
	PerformanceTier         string `json:"performanceTier"`
	PerformanceTierCategory string `json:"performanceTierCategory"`
	ActivityURL             string `json:"activityUrl"`
	CanEdit                 bool   `json:"canEdit"`
	CanCreateDatasets       bool   `json:"canCreateDatasets"`
	ActivePane              string `json:"activePane"`
	TotalDatasets           int    `json:"totalDatasets"`
	TotalOrganizations      int    `json:"totalOrganizations"`
	CompetitionsSummary     struct {
		Tier              string `json:"tier"`
		TotalResults      int    `json:"totalResults"`
		RankOutOf         int    `json:"rankOutOf"`
		RankCurrent       int    `json:"rankCurrent"`
		RankHighest       int    `json:"rankHighest"`
		TotalGoldMedals   int    `json:"totalGoldMedals"`
		TotalSilverMedals int    `json:"totalSilverMedals"`
		TotalBronzeMedals int    `json:"totalBronzeMedals"`
		SummaryType       string `json:"summaryType"`
	} `json:"competitionsSummary"`
	ScriptsSummary struct {
		Tier              string `json:"tier"`
		TotalResults      int    `json:"totalResults"`
		RankOutOf         int    `json:"rankOutOf"`
		RankCurrent       int    `json:"rankCurrent"`
		RankHighest       int    `json:"rankHighest"`
		TotalGoldMedals   int    `json:"totalGoldMedals"`
		TotalSilverMedals int    `json:"totalSilverMedals"`
		TotalBronzeMedals int    `json:"totalBronzeMedals"`
		SummaryType       string `json:"summaryType"`
	} `json:"scriptsSummary"`
	DiscussionsSummary struct {
		Tier              string `json:"tier"`
		TotalResults      int    `json:"totalResults"`
		RankOutOf         int    `json:"rankOutOf"`
		RankCurrent       int    `json:"rankCurrent"`
		RankHighest       int    `json:"rankHighest"`
		TotalGoldMedals   int    `json:"totalGoldMedals"`
		TotalSilverMedals int    `json:"totalSilverMedals"`
		TotalBronzeMedals int    `json:"totalBronzeMedals"`
		SummaryType       string `json:"summaryType"`
	} `json:"discussionsSummary"`
	Followers struct {
		Type  string `json:"type"`
		Count int    `json:"count"`
	} `json:"followers"`
	Following struct {
		Type  string `json:"type"`
		Count int    `json:"count"`
	} `json:"following"`
}

type KaggleErr struct {
	ProxyIP string `json:"proxy_ip"`
	IsError bool   `json:"is_error"`
}
