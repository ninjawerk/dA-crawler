/**
 * Created by BeastSanchez on 12/2/2017
 */

package models

import "github.com/jinzhu/gorm"

type Artwork struct {
	gorm.Model
	Title     string
	Artist    string
	Url       string
	ImageUrl  string
	FavCount  int
	ArtistUrl string
}
