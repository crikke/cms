package loader

// type mockRepo struct {
// 	content []repository.ContentData
// }

// func newMockRepo() mockRepo {

// 	repo := mockRepo{}

// 	cd := repository.ContentData{
// 		ID:               uuid.UUID{},
// 		PublishedVersion: 0,
// 		Data: map[int]repository.ContentVersion{
// 			0: {
// 				Name: map[string]string{
// 					"sv": "foo",
// 				},
// 				URLSegment: map[string]string{
// 					"sv": "foo",
// 				},
// 				Properties: []repository.ContentProperty{
// 					{
// 						ID:        uuid.UUID{},
// 						Name:      "prop",
// 						Type:      "text",
// 						Localized: false,
// 						Value: map[string]interface{}{
// 							"sv": "bar",
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// 	repo.content = []repository.ContentData{cd}

// 	return repo
// }
// func (m mockRepo) GetContent(ctx context.Context, contentReference domain.ContentReference) (repository.ContentData, error) {
// 	return m.content[0], nil
// }

// func (m mockRepo) GetChildren(ctx context.Context, contentReference domain.ContentReference) ([]repository.ContentData, error) {
// 	return nil, nil
// }

// func (m mockRepo) LoadSiteConfiguration(ctx context.Context) (*domain.SiteConfiguration, error) {
// 	return nil, nil
// }
