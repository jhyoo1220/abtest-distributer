package dbs

const (
	TESTLIST_KEY  = "testlist"
	TESTCOUNT_KEY = "testcount"

	TEST_KEY           = "test"
	TEST_NUM_USERS_KEY = "test_num_users"
)

func GetTestlistKey() string {
	return TESTLIST_KEY
}

func GetTestCountKey() string {
	return TESTCOUNT_KEY
}

func GetTestKey(name string) string {
	var b bytes.Buffer

	b.WriteString(TEST_KEY)
	b.WriteString(":")
	b.WriteString(name)

	return b.String()
}

func GetTestNumUsersKey(testName string, groupName string) string {
	var b bytes.Buffer

	b.WriteString(TEST_NUM_USERS_KEY)
	b.WriteString(":")
	b.WriteString(testName)
	b.WriteString(":")
	b.WriteString(groupName)

	return b.String()
}
