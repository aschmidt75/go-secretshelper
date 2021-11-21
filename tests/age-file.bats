#!/usr/bin/env bats

@test "invoke cli - version" {
    run ../dist/go-secretshelper version
    [ "$status" -eq 0 ]
}

@test "invoke cli - run file - nonexisting" {
    run ../dist/go-secretshelper run -c nonex
    [ "$status" -ne 0 ]
}

@test "invoke cli - run file" {
    run ../dist/go-secretshelper run -c ./fixtures/fixture-2.yaml
    [ "$status" -eq 0 ]
    [ -f ./go-secrethelper-test.dat ]
    rm ./go-secrethelper-test.dat
}

@test "invoke cli - run file w/ environment" {
    VAULT_NAME=kv run ../dist/go-secretshelper -e run -c ./fixtures/fixture-3.yaml
    [ "$status" -eq 0 ]
    [ -f ./go-secrethelper-test3.dat ]
    rm ./go-secrethelper-test3.dat
}
