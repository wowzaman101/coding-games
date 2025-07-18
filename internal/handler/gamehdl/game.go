package gamehdl

import (
	"math/rand/v2"

	"github.com/gofiber/fiber/v3"
)

type handler struct {
}

type Handler interface {
	Game(ctx fiber.Ctx) error
}

func New() Handler {
	return &handler{}
}

type Card struct {
	Number int    // 1=A, 2–9, 10, 11=J, 12=Q, 13=K
	Suit   string // hearts, diamonds, clubs, spades (unused for value)
}
type Request struct {
	PlayHands [][]Card `json:"playHands"`
	GameType  int      `json:"gameType"` // 1=game1, 2=game2
}

type Response struct {
	Data []string `json:"response"`
}

func (h *handler) Game(ctx fiber.Ctx) error {
	var req Request
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var decisions []string

	switch req.GameType {
	case 1:
		decisions = game1(req.PlayHands)
	case 2:
		decisions = game2(req.PlayHands)
	default:
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid game type",
		})
	}

	// Example handler logic
	// This could be replaced with actual game logic
	return ctx.JSON(fiber.Map{
		"response": decisions,
	})
}

func game1(hands [][]Card) []string {
	decisions := make([]string, 0)
	for _, hand := range hands {
		decisions = append(decisions, decide(hand))
	}
	return decisions
}

func newDeck() []Card {
	deck := make([]Card, 0, 36)
	suits := []string{"hearts", "diamonds", "clubs", "spades"}
	for _, suit := range suits {
		for i := 1; i <= 9; i++ {
			deck = append(deck, Card{Number: i, Suit: suit})
		}
	}
	return deck
}

func removeCard(deck []Card, card Card) []Card {
	for i, c := range deck {
		if c.Number == card.Number && c.Suit == card.Suit {
			return append(deck[:i], deck[i+1:]...)
		}
	}
	return deck // return unchanged if not found
}

func game2(hands [][]Card) []string {
	deck := newDeck()
	for _, hand := range hands {
		for _, card := range hand {
			deck = removeCard(deck, card) // remove cards from the deck
		}
	}

	decisions := make([]string, 0)
	for _, hand := range hands {
		decisions = append(decisions, decide2(hand, deck))
	}
	return decisions
}

// calculateValue returns the baccarat point total of a hand.
func calculateValue(hand []Card) int {
	sum := 0
	for _, c := range hand {
		v := c.Number
		if v >= 10 { // 10, J, Q, K all count as zero
			v = 0
		}
		sum += v
	}
	return sum % 10
}

// decide applies the baccarat rule: draw on 0–5, stand on 6–9.
func decide(hand []Card) string {
	total := calculateValue(hand)
	if total < 5 || (total == 5 && rand.Float64() < 0.3) {
		return "hit"
	}
	return "stand"
}

func decide2(hand []Card, deck []Card) string {
	total := calculateValue(hand)
	probability := 0.0
	for _, card := range deck {
		if (card.Number+total)%10 > total {
			probability += 1.0 / float64(len(deck))
		}
	}
	if probability > 0.4 {
		return "hit"
	}
	return "stand"
}
