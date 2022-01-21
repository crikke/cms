package node

import (
	"context"
	"strings"
	"time"

	"github.com/crikke/cms/pkg/locale"
	"github.com/google/uuid"
	"golang.org/x/text/language"
)

type key int

var nodeKey key

type Property struct {
	ID        uuid.UUID
	Name      string
	Type      string
	Localized bool
	// if a property isnt localized, the value will be stored in the configured default language.
	// this makes it possible to later on enable localization.
	Value map[language.Tag]interface{}
}

type Node struct {
	// A node are required to have a localized URLSegment for each configured locale.
	ID       uuid.UUID
	ParentID uuid.UUID
	Content  map[int]ContentData
	Version  int
	Created  time.Time
	Updated  time.Time
}

type ContentData struct {
	Version    int
	Name       string
	URLSegment map[string]string
	Properties []Property
}

// Checks if a node matches a urlsegment
func (n Node) Match(ctx context.Context, remaining []string) (match bool, segments []string) {
	segments = remaining
	lang := locale.FromContext(ctx)

	segment := remaining[0]
	nodeSegment, exist := n.Content[n.Version].URLSegment[lang.String()]
	// this should never fail because if URLSegment is not set explicitly, the localized name will be used as URLSegment
	if !exist {
		match = false
		return
	}

	match = strings.EqualFold(segment, nodeSegment)

	if match {
		// pop matched segment
		segments = segments[1:]
	}

	return
}

func WithNode(ctx context.Context, node Node) context.Context {
	return context.WithValue(ctx, nodeKey, node)
}

func RoutedNode(ctx context.Context) Node {
	return ctx.Value(nodeKey).(Node)
}

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
