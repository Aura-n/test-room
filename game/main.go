package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Player の構造体
type Player struct {
	ID          int
	Name        string
	Hp          int
	DefaultHp   int
	Atk         int
	DefaultAtk  int
	FirstAttack bool
	Weapon
}

func newPlayer(ID int, hp int, atk int) Player {
	p := new(Player)
	p.ID = ID
	p.Hp = hp
	p.DefaultHp = hp
	p.Atk = atk
	p.DefaultAtk = atk
	p.FirstAttack = false
	return *p
}

// Weapon の構造体
type Weapon struct {
	ID     int
	Name   string
	Atk    int
	Effect func(player *Player) Player
	Desc   string
}

func newWeapon(ID int, name string, atk int, effect func(player *Player) Player, desc string) Weapon {
	w := new(Weapon)
	w.ID = ID
	w.Name = name
	w.Atk = atk
	w.Effect = effect
	w.Desc = desc
	return *w
}

func main() {
	var baseHp, baseAtk int

	// Weapon 初期化
	weapons := []Weapon{
		newWeapon(1, "Sword", 3, func(player *Player) Player { player.Hp += 5; return *player }, "Raise a shield!(HP+5)"),
		newWeapon(2, "Dagger", 3, func(player *Player) Player { player.FirstAttack = true; return *player }, "Increased quickness!(Can First Attack)"),
		newWeapon(3, "Wand", 5, func(player *Player) Player { player.Hp -= 3; return *player }, "Dody weaken...(HP-3)"),
	}

	rand.Seed(time.Now().Unix())
	baseHp = 20
	baseAtk = 1
	// Player 作成
	player1 := newPlayer(1, baseHp+rand.Intn(20), baseAtk+rand.Intn(3))
	player2 := newPlayer(2, baseHp+rand.Intn(20), baseAtk+rand.Intn(3))

	// Game start
	enterName(&player1)
	enterName(&player2)
	showPlayerStatus(player1)
	showPlayerStatus(player2)
	selectWeapon(weapons, &player1)
	selectWeapon(weapons, &player2)
	applyEffect(&player1)
	applyEffect(&player2)
	showPlayerStatus(player1)
	showPlayerStatus(player2)
	fmt.Print("Almost start the game.")
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second * 1)
		fmt.Print(".")
	}
	fmt.Println()
	startBattle(player1, player2)
}

func enterName(player *Player) {
	// 名前入力
	fmt.Printf("Player %d Enter your name. >", player.ID)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	player.Name = scanner.Text()
	fmt.Printf("Player %d name is %s.\n", player.ID, player.Name)
}

func selectWeapon(weapons []Weapon, player *Player) {
	var selectedWeaponID int

SELECTWEAPON:
	// 武器表示
	for _, weapon := range weapons {
		fmt.Printf("%d. Name:%s ATK:%d DESC:%s\n", weapon.ID, weapon.Name, weapon.Atk, weapon.Desc)
	}
	fmt.Printf("%s select your weapon. (Enter the weapon number) >", player.Name)
	fmt.Scan(&selectedWeaponID)

	// 武器確認
	for _, weapon := range weapons {
		if weapon.ID == selectedWeaponID {
			player.Weapon = weapon
			player.Atk = player.Atk + player.Weapon.Atk
			fmt.Printf("You selected %s!\n", player.Weapon.Name)
			return
		}
	}
	// 入力した武器 ID が存在しない場合の処理
	fmt.Println("Your select weapon number is nothing.")
	goto SELECTWEAPON
}

func showPlayerStatus(player Player) {
	// Player 情報の表示
	fmt.Printf("Name: %s HP: %d(%d) ATK: %d(%d) Weapon: %s\n", player.Name, player.Hp, player.Hp-player.DefaultHp, player.Atk, player.Atk-player.DefaultAtk, player.Weapon.Name)
}

func startBattle(player Player, opponent Player) {

	playerHp := player.Hp
	opponentHp := opponent.Hp

	// 先制攻撃
	if player.FirstAttack == true {
		time.Sleep(time.Second * 1)
		// player の攻撃
		opponentHp -= player.Atk
		fmt.Println("First Attack!!")
		fmt.Printf("%s HP: %d ATK: %d ===> %s HP: %d ATK: %d \n", player.Name, playerHp, player.Atk, opponent.Name, opponentHp, opponent.Atk)
		if judge(player, opponent, playerHp, opponentHp) {
			return
		}
	}

	if opponent.FirstAttack == true {
		time.Sleep(time.Second * 1)
		// opponent の攻撃
		playerHp -= opponent.Atk
		fmt.Println("First Attack!!")
		fmt.Printf("%s HP: %d ATK: %d <=== %s HP: %d ATK: %d \n", player.Name, playerHp, player.Atk, opponent.Name, opponentHp, opponent.Atk)
		if judge(player, opponent, playerHp, opponentHp) {
			return
		}
	}

	for {
		time.Sleep(time.Second * 1)
		// player の攻撃
		opponentHp -= player.Atk
		fmt.Printf("%s HP: %d ATK: %d ===> %s HP: %d ATK: %d \n", player.Name, playerHp, player.Atk, opponent.Name, opponentHp, opponent.Atk)
		if judge(player, opponent, playerHp, opponentHp) {
			return
		}

		time.Sleep(time.Second * 1)
		// opponent の攻撃
		playerHp -= opponent.Atk
		fmt.Printf("%s HP: %d ATK: %d <=== %s HP: %d ATK: %d \n", player.Name, playerHp, player.Atk, opponent.Name, opponentHp, opponent.Atk)
		if judge(player, opponent, playerHp, opponentHp) {
			return
		}
	}
}

func judge(player, opponent Player, playerHp, opponentHp int) bool {

	// Battle
	if playerHp <= 0 && opponentHp <= 0 {
		// 引き分け
		fmt.Println("This game is Draw!!")
		return true
	} else if playerHp <= 0 {
		// opponent の勝利
		fmt.Printf("%s is Win!!", opponent.Name)
		return true
	} else if opponentHp <= 0 {
		// player の勝利
		fmt.Printf("%s is Win!!", player.Name)
		return true
	}
	return false
}

func applyEffect(player *Player) {
	// Effect 適用
	player.Effect(player)
}
