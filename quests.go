package main

type Quest struct {
	ID          int
	Description string
	GoalType    string
	Target      string
	Required    int
	Progress    int
	RewardXP    int
	Completed   bool
}
