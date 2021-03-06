package storage

import (
	"strings"
	"sync"
)

type ChannelStore struct {
	users    map[string]map[string][]string
	userLock sync.Mutex

	topic     map[string]map[string]string
	topicLock sync.Mutex
}

func NewChannelStore() *ChannelStore {
	return &ChannelStore{
		users: make(map[string]map[string][]string),
		topic: make(map[string]map[string]string),
	}
}

func (c *ChannelStore) GetUsers(server, channel string) []string {
	c.userLock.Lock()

	users := make([]string, len(c.users[server][channel]))
	copy(users, c.users[server][channel])

	c.userLock.Unlock()

	return users
}

func (c *ChannelStore) SetUsers(users []string, server, channel string) {
	c.userLock.Lock()

	if _, ok := c.users[server]; !ok {
		c.users[server] = make(map[string][]string)
	}

	c.users[server][channel] = users
	c.userLock.Unlock()
}

func (c *ChannelStore) AddUser(user, server, channel string) {
	c.userLock.Lock()

	if _, ok := c.users[server]; !ok {
		c.users[server] = make(map[string][]string)
	}

	if users, ok := c.users[server][channel]; ok {
		for _, nick := range users {
			if nick == user {
				c.userLock.Unlock()
				return
			}
		}

		c.users[server][channel] = append(users, user)
	} else {
		c.users[server][channel] = []string{user}
	}

	c.userLock.Unlock()
}

func (c *ChannelStore) RemoveUser(user, server, channel string) {
	c.userLock.Lock()
	c.removeUser(user, server, channel)
	c.userLock.Unlock()
}

func (c *ChannelStore) RemoveUserAll(user, server string) {
	c.userLock.Lock()

	for channel := range c.users[server] {
		c.removeUser(user, server, channel)
	}

	c.userLock.Unlock()
}

func (c *ChannelStore) RenameUser(oldNick, newNick, server string) {
	c.userLock.Lock()
	c.renameAll(server, oldNick, newNick)
	c.userLock.Unlock()
}

func (c *ChannelStore) SetMode(server, channel, user, add, remove string) {
	c.userLock.Lock()

	if strings.Contains(add, "o") {
		c.setPrefix(server, channel, user, "@")
	} else if strings.Contains(add, "v") {
		c.setPrefix(server, channel, user, "+")
	} else if strings.IndexAny(remove, "ov") > -1 {
		c.setPrefix(server, channel, user, "")
	}

	c.userLock.Unlock()
}

func (c *ChannelStore) FindUserChannels(user, server string) []string {
	var channels []string

	c.userLock.Lock()

	for channel, users := range c.users[server] {
		for _, nick := range users {
			if strings.TrimLeft(nick, "@+") == user {
				channels = append(channels, channel)
				break
			}
		}
	}

	c.userLock.Unlock()

	return channels
}

func (c *ChannelStore) GetTopic(server, channel string) string {
	c.topicLock.Lock()
	defer c.topicLock.Unlock()

	return c.topic[server][channel]
}

func (c *ChannelStore) SetTopic(topic, server, channel string) {
	c.topicLock.Lock()

	if _, ok := c.topic[server]; !ok {
		c.topic[server] = make(map[string]string)
	}

	c.topic[server][channel] = topic
	c.topicLock.Unlock()
}

func (c *ChannelStore) rename(server, channel, oldNick, newNick string) {
	for i, nick := range c.users[server][channel] {
		if strings.TrimLeft(nick, "@+") == oldNick {
			if nick[0] == '@' || nick[0] == '+' {
				newNick = nick[:1] + newNick
			}

			c.users[server][channel][i] = newNick
			return
		}
	}
}

func (c *ChannelStore) setPrefix(server, channel, user, prefix string) {
	for i, nick := range c.users[server][channel] {
		if strings.TrimLeft(nick, "@+") == user {
			c.users[server][channel][i] = prefix + user
			return
		}
	}
}

func (c *ChannelStore) renameAll(server, oldNick, newNick string) {
	for channel := range c.users[server] {
		c.rename(server, channel, oldNick, newNick)
	}
}

func (c *ChannelStore) removeUser(user, server, channel string) {
	for i, nick := range c.users[server][channel] {
		if strings.TrimLeft(nick, "@+") == user {
			users := c.users[server][channel]
			c.users[server][channel] = append(users[:i], users[i+1:]...)
			return
		}
	}
}
