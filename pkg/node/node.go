package node

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

type key int

var nodeKey key

type Property struct {
	ID        uuid.UUID
	Name      string
	Type      string
	Localized bool
	// detta kanske ska vara ett interface som har en transform metod som returnerar []byte och converterar till json
	Value interface{}
}

// det är contentloader som hanterar hämtningen av korrekt localized / version av node.
// alltså propertytransform vid hämtning från db och inte vid json marshal
// Därför ska Node (Content) inte ha några maps eller version, utan endast datan som ska skickas
// ansvaret för att Mappa DbContent till Content görs alltså av contentloader
type Node struct {
	// A node are required to have a localized URLSegment for each configured locale.
	ID         uuid.UUID
	ParentID   uuid.UUID
	URLSegment string
	Name       string
	Properties []Property
	Created    time.Time
	Updated    time.Time
}

type ContentData struct {
	Version    int
	Name       string
	Properties []Property
}

// Checks if a node matches a urlsegment
func (n Node) Match(ctx context.Context, remaining []string) (match bool, segments []string) {
	segments = remaining

	segment := remaining[0]

	nodeSegment := n.URLSegment

	match = strings.EqualFold(segment, nodeSegment)

	if match {
		// pop matched segment
		segments = segments[1:]
	}

	return
}

// func ContentResultHandler(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
// 		content := Content(r.Context())
// 		node := RoutedNode(r.Context())

// 		var properies []interface{}

// 		for _, p := range content.Properties {

// 			if prop, exist := content.Properties.Value[locale.FromContext(r.Context())]; exist {
// 				properies = append(properies, prop)
// 			} else {
// 				properies = append(properies, p.Value[config.SiteConfiguration.Languages[0]])
// 			}

// 		}

// 		response := ContentResponse{
// 			node,
// 			"",
// 			properies,
// 		}
// 	})
// }

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
