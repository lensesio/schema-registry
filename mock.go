package schemaregistry

type MockClient struct {
	SubjectsFn           func() (subjects []string, err error)
	VersionsFn           func(subject string) (versions []int, err error)
	RegisterNewSchemaFn  func(subject, schema string) (int, error)
	IsRegisteredFn       func(subject, schema string) (bool, Schema, error)
	GetSchemaByIdFn      func(id int) (string, error)
	GetSchemaBySubjectFn func(subject string, ver int) (Schema, error)
	GetLatestSchemaFn    func(subject string) (Schema, error)
}

func (c *MockClient) Subjects() (subjects []string, err error) {
	return c.SubjectsFn()
}

func (c *MockClient) Versions(subject string) (versions []int, err error) {
	return c.VersionsFn(subject)
}

func (c *MockClient) RegisterNewSchema(subject, schema string) (int, error) {
	return c.RegisterNewSchemaFn(subject, schema)
}

func (c *MockClient) IsRegistered(subject, schema string) (bool, Schema, error) {
	return c.IsRegisteredFn(subject, schema)
}

func (c *MockClient) GetSchemaById(id int) (string, error) {
	return c.GetSchemaByIdFn(id)
}

func (c *MockClient) GetSchemaBySubject(subject string, ver int) (Schema, error) {
	return c.GetSchemaBySubjectFn(subject, ver)
}

func (c *MockClient) GetLatestSchema(subject string) (Schema, error) {
	return c.GetLatestSchemaFn(subject)
}

func NewNOOPClient() Client {
	return &MockClient{
		SubjectsFn: func() (subjects []string, err error) {
			return []string{}, nil
		},
		VersionsFn: func(subject string) (versions []int, err error) {
			return []int{}, nil
		},
		RegisterNewSchemaFn: func(subject, schema string) (int, error) {
			return 0, nil
		},
		IsRegisteredFn: func(subject, schema string) (bool, Schema, error) {
			return false, Schema{}, nil
		},
		GetSchemaByIdFn: func(id int) (string, error) {
			return "", nil
		},
		GetSchemaBySubjectFn: func(subject string, ver int) (Schema, error) {
			return Schema{}, nil
		},
		GetLatestSchemaFn: func(subject string) (Schema, error) {
			return Schema{}, nil
		},
	}
}
