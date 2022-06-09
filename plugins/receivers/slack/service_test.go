package slack

// type SlackRepositoryTestSuite struct {
// 	suite.Suite
// 	service SlackRepository
// 	client  *mocks.SlackCaller
// }

// func TestSlackRepository(t *testing.T) {
// 	suite.Run(t, new(SlackRepositoryTestSuite))
// }

// func (s *SlackRepositoryTestSuite) TestGetWorkspaceChannel() {
// 	oldClientCreator := newClient
// 	mockedSlackClient := &mocks.SlackCaller{}
// 	newClient = func(string) SlackCaller {
// 		return mockedSlackClient
// 	}
// 	defer func() { newClient = oldClientCreator }()
// 	s.client = mockedSlackClient
// 	s.service = NewService()

// 	s.Run("should return joined channel list in a workspace", func() {
// 		s.client.
// 			On("GetConversationsForUser", mock.AnythingOfType("*slack.GetConversationsForUserParameters")).
// 			Run(func(args mock.Arguments) {
// 				r := args.Get(0).(*slack.GetConversationsForUserParameters)
// 				s.Equal(1000, r.Limit)
// 				s.Equal([]string{"public_channel", "private_channel"}, r.Types)
// 				s.Equal("", r.Cursor)
// 			}).
// 			Return([]slack.Channel{
// 				{GroupConversation: slack.GroupConversation{Name: "foo"}},
// 				{GroupConversation: slack.GroupConversation{Name: "bar"}}}, "nextCurr", nil).Once()
// 		s.client.
// 			On("GetConversationsForUser", mock.AnythingOfType("*slack.GetConversationsForUserParameters")).
// 			Run(func(args mock.Arguments) {
// 				r := args.Get(0).(*slack.GetConversationsForUserParameters)
// 				s.Equal(1000, r.Limit)
// 				s.Equal([]string{"public_channel", "private_channel"}, r.Types)
// 				s.Equal("nextCurr", r.Cursor)
// 			}).
// 			Return([]slack.Channel{
// 				{GroupConversation: slack.GroupConversation{Name: "baz"}}}, "", nil).Once()
// 		channels, err := s.service.GetWorkspaceChannels("test_token")
// 		s.Equal(3, len(channels))
// 		s.Equal("foo", channels[0].Name)
// 		s.Equal("bar", channels[1].Name)
// 		s.Equal("baz", channels[2].Name)
// 		s.Nil(err)
// 		s.client.AssertExpectations(s.T())
// 	})

// 	s.Run("should return error if get joined channel list fail", func() {
// 		s.client.On("GetConversationsForUser", mock.Anything).
// 			Return(nil, "", errors.New("random error")).Once()

// 		channels, err := s.service.GetWorkspaceChannels("test_token")
// 		s.Nil(channels)
// 		s.EqualError(err, "failed to fetch joined channel list: random error")
// 	})
// }

// func (s *SlackRepositoryTestSuite) TestNotify() {
// 	oldClientCreator := newClient
// 	mockedSlackClient := &mocks.SlackCaller{}
// 	newClient = func(string) SlackCaller {
// 		return mockedSlackClient
// 	}
// 	defer func() { newClient = oldClientCreator }()
// 	s.client = mockedSlackClient
// 	s.service = NewService()

// 	s.Run("should notify user identified by their email", func() {
// 		mockedSlackClient.On("GetUserByEmail", "foo@odpf.io").
// 			Return(&slack.User{ID: "U20"}, nil).Once()
// 		mockedSlackClient.On("SendMessage", "U20",
// 			mock.AnythingOfType("slack.MsgOption"), mock.AnythingOfType("slack.MsgOption")).Return("", "", "", nil).Once()
// 		dummyMessage := &SlackMessage{
// 			ReceiverName: "foo@odpf.io",
// 			ReceiverType: "user",
// 			Message:      "random text",
// 			Token:        "foo_bar",
// 		}
// 		res, err := s.service.Notify(dummyMessage)
// 		s.Nil(err)
// 		s.True(res.OK)
// 	})

// 	s.Run("should return error if notifying user fails", func() {
// 		mockedSlackClient.On("GetUserByEmail", "foo@odpf.io").
// 			Return(&slack.User{ID: "U20"}, nil).Once()
// 		mockedSlackClient.On("SendMessage", "U20",
// 			mock.AnythingOfType("slack.MsgOption"),
// 			mock.AnythingOfType("slack.MsgOption"),
// 		).Return("", "", "", errors.New("random error")).Once()

// 		dummyMessage := &SlackMessage{
// 			ReceiverName: "foo@odpf.io",
// 			ReceiverType: "user",
// 			Message:      "random text",
// 			Token:        "foo_bar",
// 		}
// 		res, err := s.service.Notify(dummyMessage)
// 		s.EqualError(err, "could not send notification: failed to send message to foo@odpf.io: random error")
// 		s.False(res.OK)
// 	})

// 	s.Run("should return error if user lookup by email fails", func() {
// 		mockedSlackClient.On("GetUserByEmail", "foo@odpf.io").
// 			Return(nil, errors.New("users_not_found")).Once()

// 		dummyMessage := &SlackMessage{
// 			ReceiverName: "foo@odpf.io",
// 			ReceiverType: "user",
// 			Message:      "random text",
// 			Token:        "foo_bar",
// 		}
// 		res, err := s.service.Notify(dummyMessage)
// 		s.EqualError(err, "could not send notification: failed to get id for foo@odpf.io: users_not_found")
// 		s.False(res.OK)
// 	})

// 	s.Run("should return error if user lookup by email returns any error", func() {
// 		mockedSlackClient.On("GetUserByEmail", "foo@odpf.io").
// 			Return(nil, errors.New("random error")).Once()

// 		dummyMessage := &SlackMessage{
// 			ReceiverName: "foo@odpf.io",
// 			ReceiverType: "user",
// 			Message:      "random text",
// 			Token:        "foo_bar",
// 		}
// 		res, err := s.service.Notify(dummyMessage)
// 		s.EqualError(err, "could not send notification: random error")
// 		s.False(res.OK)
// 	})

// 	s.Run("should notify if part of the channel", func() {
// 		mockedSlackClient.On("GetConversationsForUser", mock.Anything).Return([]slack.Channel{
// 			{GroupConversation: slack.GroupConversation{
// 				Name:         "foo",
// 				Conversation: slack.Conversation{ID: "C01"}},
// 			}, {GroupConversation: slack.GroupConversation{
// 				Name:         "bar",
// 				Conversation: slack.Conversation{ID: "C02"}},
// 			}}, "", nil).Once()

// 		mockedSlackClient.On("SendMessage", "C01",
// 			mock.AnythingOfType("slack.MsgOption"),
// 			mock.AnythingOfType("slack.MsgOption"),
// 		).Return("", "", "", nil).Once()

// 		dummyMessage := &SlackMessage{
// 			ReceiverName: "foo",
// 			ReceiverType: "channel",
// 			Message:      "random text",
// 			Token:        "foo_bar",
// 		}
// 		res, err := s.service.Notify(dummyMessage)
// 		s.Nil(err)
// 		s.True(res.OK)
// 		mockedSlackClient.AssertExpectations(s.T())
// 	})

// 	s.Run("should return error if not part of the channel", func() {
// 		mockedSlackClient.On("GetConversationsForUser", mock.Anything).Return([]slack.Channel{
// 			{GroupConversation: slack.GroupConversation{
// 				Name:         "foo",
// 				Conversation: slack.Conversation{ID: "C01"}},
// 			}, {GroupConversation: slack.GroupConversation{
// 				Name:         "bar",
// 				Conversation: slack.Conversation{ID: "C02"}},
// 			}}, "", nil).Once()

// 		mockedSlackClient.On("SendMessage", "C01",
// 			mock.AnythingOfType("slack.MsgOption")).Return("", "", "", nil).Once()

// 		dummyMessage := &SlackMessage{
// 			ReceiverName: "baz",
// 			ReceiverType: "channel",
// 			Message:      "random text",
// 			Token:        "foo_bar",
// 		}
// 		res, err := s.service.Notify(dummyMessage)
// 		s.EqualError(err, "could not send notification: app is not part of the channel baz")
// 		s.False(res.OK)
// 	})

// 	s.Run("should return error failed to fetch joined channels list", func() {
// 		mockedSlackClient.On("GetConversationsForUser", mock.Anything).
// 			Return(nil, "", errors.New("random error")).Once()
// 		mockedSlackClient.On("SendMessage", "C01",
// 			mock.AnythingOfType("slack.MsgOption")).Return("", "", "", nil).Once()

// 		dummyMessage := &SlackMessage{
// 			ReceiverName: "baz",
// 			ReceiverType: "channel",
// 			Message:      "random text",
// 			Token:        "foo_bar",
// 		}
// 		res, err := s.service.Notify(dummyMessage)
// 		s.EqualError(err, "could not send notification: failed to fetch joined channel list: random error")
// 		s.False(res.OK)
// 	})
// }
