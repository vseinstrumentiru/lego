package di

// Provider returns provider - function for registering provider and configurations - functions for configure it
type Module func() (provider interface{}, configurations []interface{})
