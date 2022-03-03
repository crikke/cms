package command

/*
example:

POST /contentdefinition/propertydefinition
{
	cid: uuid,
	name: name
	description: desc,
	type: text,
}
*/

/*
PUT /contentdefinition/{cid}/propertydefinition/{pid}
{
	required: true,
	regex: "^(foo)",
	unique: true,
	localized: true,
	boundary: { // coordinates - polygon in which a coord needs to be to be valid
		x,y
		x,y
		x,y
	}
}
// */

// type CreatePropertyDefinition struct {
// 	ContentDefinitionID uuid.UUID
// 	Name                string
// 	Description         string
// 	Type                string
// }

// type CreatePropertyDefinitionHandler struct {
// 	Repo contentdefinition.ContentDefinitionRepository
// }

// func (h CreatePropertyDefinitionHandler) Handle(ctx context.Context, cmd CreatePropertyDefinition) (uuid.UUID, error) {

// 	if cmd.ContentDefinitionID == (uuid.UUID{}) {
// 		return uuid.UUID{}, errors.New("empty contentdefinition id")
// 	}

// 	pd, err := contentdefinition.NewPropertyDefinition(cmd.Name, cmd.Description, cmd.Type)

// 	if err != nil {
// 		return uuid.UUID{}, err
// 	}

// 	pd.ID = uuid.New()
// 	err = h.Repo.UpdateContentDefinition(
// 		ctx,
// 		cmd.ContentDefinitionID,
// 		func(ctx context.Context, cd *contentdefinition.ContentDefinition) (*contentdefinition.ContentDefinition, error) {

// 			cd.Propertydefinitions = append(cd.Propertydefinitions, pd)
// 			return cd, nil
// 		})

// 	if err != nil {
// 		return uuid.UUID{}, err
// 	}

// 	return pd.ID, nil
// }

// type UpdatePropertyDefinition struct {
// 	ContentDefinitionID  uuid.UUID
// 	PropertyDefinitionID uuid.UUID
// 	Name                 *string
// 	Description          *string
// 	Localized            *bool
// }

// type UpdatePropertyDefinitionHandler struct {
// 	repo contentdefinition.ContentDefinitionRepository
// }

// func (h UpdatePropertyDefinitionHandler) Handle(ctx context.Context, cmd UpdatePropertyDefinition) error {

// 	if cmd.ContentDefinitionID == (uuid.UUID{}) {
// 		return errors.New("empty contentdefinition id")
// 	}

// 	if cmd.PropertyDefinitionID == (uuid.UUID{}) {
// 		return errors.New("empty propertydefinition id")
// 	}

// 	return h.repo.UpdateContentDefinition(
// 		ctx, cmd.ContentDefinitionID,
// 		func(ctx context.Context, cd *contentdefinition.ContentDefinition) (*contentdefinition.ContentDefinition, error) {

// 			pd := contentdefinition.PropertyDefinition{}
// 			idx := -1
// 			for i, p := range cd.Propertydefinitions {
// 				if p.ID == cmd.PropertyDefinitionID {
// 					pd = p
// 					idx = i
// 					break
// 				}
// 			}
// 			if idx == -1 {
// 				return nil, errors.New("propertydefinition not found")
// 			}

// 			if cmd.Description != nil {
// 				pd.Description = *cmd.Description
// 			}

// 			if cmd.Name != nil && *cmd.Name != "" {
// 				pd.Name = *cmd.Name
// 			}

// 			if cmd.Localized != nil {
// 				pd.Localized = *cmd.Localized
// 			}

// 			cd.Propertydefinitions[idx] = pd
// 			return cd, nil
// 		})
// }

// type DeletePropertyDefinition struct {
// 	ContentDefinitionID  uuid.UUID
// 	PropertyDefinitionID uuid.UUID
// }

// type DeletePropertyDefinitionHandler struct {
// 	repo contentdefinition.ContentDefinitionRepository
// }

// func (h DeletePropertyDefinitionHandler) Handle(ctx context.Context, cmd DeletePropertyDefinition) error {

// 	if cmd.ContentDefinitionID == (uuid.UUID{}) {
// 		return errors.New("empty contentdefinition id")
// 	}

// 	if cmd.PropertyDefinitionID == (uuid.UUID{}) {
// 		return errors.New("empty propertydefinition id")
// 	}

// 	return h.repo.UpdateContentDefinition(
// 		ctx, cmd.ContentDefinitionID,
// 		func(ctx context.Context, cd *contentdefinition.ContentDefinition) (*contentdefinition.ContentDefinition, error) {

// 			idx := -1
// 			for i, p := range cd.Propertydefinitions {
// 				if p.ID == cmd.PropertyDefinitionID {
// 					idx = i
// 					break
// 				}
// 			}
// 			if idx == -1 {
// 				return nil, errors.New("propertydefinition not found")
// 			}

// 			arr := cd.Propertydefinitions
// 			arr[idx] = arr[len(arr)-1]
// 			arr = arr[:len(arr)-1]

// 			cd.Propertydefinitions = arr
// 			return cd, nil
// 		})
// }

// type UpdateValidator struct {
// 	ContentDefinitionID  uuid.UUID
// 	PropertyDefinitionID uuid.UUID
// 	ValidatorName        string
// 	Value                interface{}
// }

// type UpdateValidatorHandler struct {
// 	Repo contentdefinition.ContentDefinitionRepository
// }

// func (h UpdateValidatorHandler) Handle(ctx context.Context, cmd UpdateValidator) error {

// 	v, err := validator.Parse(cmd.ValidatorName, cmd.Value)

// 	if err != nil {
// 		return err
// 	}

// 	return h.Repo.UpdateContentDefinition(ctx,
// 		cmd.ContentDefinitionID,
// 		func(ctx context.Context, cd *contentdefinition.ContentDefinition) (*contentdefinition.ContentDefinition, error) {

// 			idx := 0
// 			pd := contentdefinition.PropertyDefinition{}
// 			for i, p := range cd.Propertydefinitions {

// 				if p.ID == cmd.PropertyDefinitionID {
// 					pd = p
// 					idx = i
// 					break
// 				}
// 			}

// 			if pd.ID == (uuid.UUID{}) {
// 				return nil, errors.New("propertydefinition not found")
// 			}

// 			if pd.Validators == nil {
// 				pd.Validators = make(map[string]interface{})
// 			}

// 			pd.Validators[cmd.ValidatorName] = v
// 			cd.Propertydefinitions[idx] = pd
// 			return cd, nil
// 		})

// 	// return h.Repo.UpdatePropertyDefinition(
// 	// 	ctx,
// 	// 	cmd.ContentDefinitionID,
// 	// 	cmd.PropertyDefinitionID,
// 	// 	func(ctx context.Context, pd *contentdefinition.PropertyDefinition) (*contentdefinition.PropertyDefinition, error) {
// 	// 		if pd.Validators == nil {
// 	// 			pd.Validators = make(map[string]interface{})
// 	// 		}

// 	// 		pd.Validators[cmd.ValidatorName] = v
// 	// 		return pd, nil
// 	// 	})
// }
