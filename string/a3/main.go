package main

import (
	"fmt"
	"strings"
)

func main() {
	aa := filterSensitiveWord("Zenly is a map that lets you what your friends and family are up to. Millions of people around the world have shown that Zenly brings you to spend more time with the people who matter most and moves you a little closer to them even when you can&#39;t. Use Zenly to:\n\n- See exactly where your friends are and what they’re up to\n- Know when your friends run out of battery when you can’t reach them\n- Bump phones with friends to let others know you’re hanging out\n- Know when they’re at home, at school, at work and when they’re on their way to any of those places\n- Get to your friends through your favorite transportation apps without asking for the address\n- Get notified when your friends travel to different countries\n- Blast your friends with your favorite emojis (hint: sound on)\n- Know when friends are together with fire on the map and create a group chat in a tap\n- Blur or freeze yourself to get off the grid for as long as you need\n\nYou don’t need hundreds of friends on Zenly - just a couple to get started. There are a number of ways to add friends - though your friend will need to accept your request in order for you to start seeing them on your map. \n\nZenly is free and without ads. We’ve also spent nine years developing a very special algorithm that allows for continuous sharing without draining your battery.")
	fmt.Println(aa)
}

// 过滤敏感词
func filterSensitiveWord(src string) bool {
	sensitiveWord := []string{"Putin", "Xi Jinping", "communism", "Mao Zedong", "Shit", "blockhead", "Racist", "Nazi", "oppositionist", "Snout", "Muzzle", "Fucke", "stupid kakhah", "Bitch", "fuck", "fuck you", "Criminals and terrorists", "Muslims and terrorists", "Anti-Part", "Anti-Communist", "Smear China", "Slander the country", "Heroin", "pornography", "prostitute", "Sell oneself", "Pervert", "Asshole", "Yousuck", "kick ass", "bastard", "stupid jerk", "dick", "stupid idlot", "freak", "whore", "asshole", "Damn you", "fuck you", "Nerd", "bitch", "son of bitch", "suck for you SB", "Playing with fire ", "Pervert", "stupid", "idiot", "go to hell", "Shut up", "Bullshit", "God damn it", "SOB", "Drug dealing", "Dark we"}

	for _, v := range sensitiveWord {
		any := strings.Contains(src, v)
		if any == true {
			fmt.Println(v)
			return false
		}
	}
	return true
}
