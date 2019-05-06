package kv

const ErrNotExist = KVError("key does not exist")
const ErrNoMatch = KVError("no keys match")

type KVError string

func (e KVError) Error() string { return string(e) }
