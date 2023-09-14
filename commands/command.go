package commands

type Command interface {
  ValidateParams() error
  Call() error
}
