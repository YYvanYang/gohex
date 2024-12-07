package vo

import (
	"net/url"
	"strings"
	"time"
)

type UserProfile struct {
	name      string
	bio       string
	avatar    string
	location  string
	website   string
	updatedAt time.Time
}

func NewUserProfile(name, bio string) (UserProfile, error) {
	profile := UserProfile{
		name:      strings.TrimSpace(name),
		bio:       strings.TrimSpace(bio),
		updatedAt: time.Now(),
	}

	if err := profile.validate(); err != nil {
		return UserProfile{}, err
	}

	return profile, nil
}

func (p UserProfile) Name() string {
	return p.name
}

func (p UserProfile) Bio() string {
	return p.bio
}

func (p UserProfile) Avatar() string {
	return p.avatar
}

func (p UserProfile) Location() string {
	return p.location
}

func (p UserProfile) Website() string {
	return p.website
}

func (p UserProfile) UpdatedAt() time.Time {
	return p.updatedAt
}

func (p UserProfile) WithAvatar(avatar string) (UserProfile, error) {
	if avatar != "" {
		if err := validateURL(avatar); err != nil {
			return p, ErrInvalidAvatarURL
		}
	}
	p.avatar = avatar
	p.updatedAt = time.Now()
	return p, nil
}

func (p UserProfile) WithLocation(location string) UserProfile {
	p.location = strings.TrimSpace(location)
	p.updatedAt = time.Now()
	return p
}

func (p UserProfile) WithWebsite(website string) (UserProfile, error) {
	if website != "" {
		if err := validateURL(website); err != nil {
			return p, ErrInvalidWebsiteURL
		}
	}
	p.website = website
	p.updatedAt = time.Now()
	return p, nil
}

func (p UserProfile) IsEmpty() bool {
	return p.name == ""
}

func (p UserProfile) validate() error {
	if p.name == "" {
		return ErrInvalidName
	}
	if len(p.name) > 100 {
		return ErrNameTooLong
	}
	if len(p.bio) > 500 {
		return ErrBioTooLong
	}
	return nil
}

func validateURL(rawURL string) error {
	_, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return err
	}
	return nil
} 