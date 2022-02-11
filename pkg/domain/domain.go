package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/language"
)

type User struct {
	ID   uuid.UUID
	Name string
}

func (u User) GetID() uuid.UUID {
	return u.ID
}

func (u User) GetName() string {
	return u.Name
}

type SiteConfiguration struct {
	Languages []language.Tag
	RootPage  uuid.UUID
}

type Property struct {
	ID        uuid.UUID
	Name      string
	Type      string
	Localized bool
	Value     interface{}
}

// A contentreference is a reference to a piece of content that is versioned & has a locale
type ContentReference struct {
	ID      uuid.UUID
	Version int
	Locale  *language.Tag
}

func (c ContentReference) String() string {
	l := language.Tag{}

	if c.Locale != nil {
		l = *c.Locale
	}
	return fmt.Sprintf("%s-%d-%s", c.ID, c.Version, l)
}

// det är contentloader som hanterar hämtningen av korrekt localized / version av node.
// alltså propertytransform vid hämtning från db och inte vid json marshal
// Därför ska Content (Content) inte ha några maps eller version, utan endast datan som ska skickas
// ansvaret för att Mappa DbContent till Content görs alltså av contentloader
type Content struct {
	// A node are required to have a localized URLSegment for each configured locale.
	ID         ContentReference
	ParentID   uuid.UUID
	URLSegment string
	Name       string
	Properties []Property
	Created    time.Time
	Updated    time.Time
}

type ContentDefinition struct {
	ID      uuid.UUID
	Name    string
	Created time.Time
	Updated time.Time

	Properties []PropertyDefinition
}

type PropertyDefinition struct {
}

// Checks if a node matches a urlsegment

/* json data structure
this is probably how the page will look in the database

versioning:

todo for now. When a new version is added, just make a copy of all properties each with the new version.
Later could do something like a version could point to a older version. (something like symlink, if that makes sense).


{
	id: uuid.UUID,

	created: "2020-01-01",
	publishedVersion: 2,
	status: "published"  // draft, published, archived
	parentId: uuid | null,
	data: {
		1: {
			urlSegment: {
				sv: "/hej-world",
				en: "/hello-world"
			},
			name: {
				sv: "Hejsan!",
				en: "Hello!",
			},
			properties: [
				{
					id:uuid.UUID,
					name: "header",
					type: "text",
					value: {
						sv:"Hejsan wärlden",
						en:"Hello World!",
					}
					localized: true,
				},
				{
					id:uuid.UUID,
					name: "header",
					type: "text",
					value: {
						sv:"Hejsan wärlden",
						en:"Hello World!",
					}
					localized: true,
				},
			]
		},
		2: {
			urlSegment: {
				sv: "/",
				en: "/"
			},
			name: {
				sv: "Hejsan!",
				en: "Hello!",
			},
			properties: [
				{
					id:uuid.UUID,
					name: "header",
					type: "text",
					value: {
						sv:"Hejsan wärlden",
						en:"Hello World!",
					}
					localized: true,
				},
				{
					id:uuid.UUID,
					name: "header",
					type: "text",
					value: {
						sv:"Hejsan wärlden",
						en:"Hello World!",
					}
					localized: true,
				},
			]
		},
	}
}

*/

// {"_id":{"$oid":"61fa9eb92cab5a2a027c1e95"},"status":"published","publishedVersion":0,"parentID":"","data":{"0":{"urlSegment":{"sv":"/hej-world","en":"/hello-world"},"name":{"sv":"Hejsan!","en":"Hello!"},"properties":[{"id":"uuid","name":"header","type":"text","value":{"sv":"Hejsan Värden!","en":"Hello World!"},"localized":true}]}}}
