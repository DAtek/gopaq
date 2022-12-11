pkgs := "."
coverfile := ".coverage"

test *options:
    go test {{ options }} {{ pkgs }}


test-cover *options:
    go test {{ options }} -coverprofile {{ coverfile }} {{ pkgs }}


show-coverage:
    go tool cover -html={{ coverfile }}


test-and-show-covarage: test-cover show-coverage
