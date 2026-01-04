package keyboard

import (
	// "fmt"
	"time"
	"unicode"

	"github.com/AllMeatball/u64-remote/server"
)

const (
	KEYBOARD_BUFFER_OFS = 631
	KEYBOARD_QUEUED_KEYS_OFS = 198
	KEYBOARD_SHFLAG_OFS = 653

	KEYBOARD_MAX_COUNT = 10
)

const (
	KEY_DEL_INS = 20
	KEY_RETURN  = 13

	KEY_CRSR_UP    = 145
	KEY_CRSR_DOWN  = 17
	KEY_CRSR_LEFT  = 157
	KEY_CRSR_RIGHT = 29
)

const (
	SHFLAG_SHIFT     = 1
	SHFLAG_COMMODORE = 2
	SHFLAG_CONTROL   = 4
)

func StopTyping(server server.U64Server) error {
	return server.PokeMemory(KEYBOARD_QUEUED_KEYS_OFS, []byte{0})
}

func GetQueuedKeys(server server.U64Server) (byte, error) {
	key_queue_data, err := server.PeekMemory(KEYBOARD_QUEUED_KEYS_OFS, 1)
	if err != nil { return 0, err }

	return key_queue_data[0], nil
}

func GetBufferSize(server server.U64Server) (byte, error) {
	key_queue_data, err := server.PeekMemory(KEYBOARD_QUEUED_KEYS_OFS, 1)
	keys_in_queue := key_queue_data[0]

	if err != nil {
		return keys_in_queue, err
	}


	return keys_in_queue, nil
}


func IsBufferFull(server server.U64Server) (bool, error) {
	keys_in_queue, err := GetBufferSize(server)
	if err != nil {
		return false, err
	}

	return keys_in_queue >= KEYBOARD_MAX_COUNT, nil
}

/*
 * Sets shift flag bits.
 * Usage: keyboard.SetShiftFlag(SHFLAG_SHIFT, true)
 *
 * NOTE: Currently this seems to not work (likely due to the key being raised after being set)
 */
func SetShiftFlag(server server.U64Server, mask byte, is_down bool) error {
	cur_shift_flag_bytes, err := server.PeekMemory(KEYBOARD_SHFLAG_OFS, 1)
	if err != nil { return err }

	cur_shift_flag := cur_shift_flag_bytes[0]

	if is_down {
		 cur_shift_flag |= mask
	}

	return server.PokeMemory(KEYBOARD_SHFLAG_OFS, []byte{cur_shift_flag})
}

func TypeBytes(server server.U64Server, keys []byte) error {
	slice_count  := len(keys) / 10
	extra_keys := len(keys) % 10

	keyboard_slices := [][]byte{}

	slice_end := 0
	for i := range slice_count {
		slice_start := 10 * i
		slice_end    = 10 * (i + 1)

		keyboard_slice := keys[slice_start:slice_end]
		keyboard_slices = append(keyboard_slices, keyboard_slice)
	}


	// Add rest of keys
	keyboard_slices = append(keyboard_slices, keys[slice_end:slice_end + extra_keys])

	// fmt.Println(keyboard_slices)

	for _, slice := range keyboard_slices {
		err := server.PokeMemory(KEYBOARD_BUFFER_OFS, slice)
		if err != nil { return err }

		err  = server.PokeMemory(KEYBOARD_QUEUED_KEYS_OFS, []byte{byte(len(slice))})
		if err != nil { return err }

		key_queue_full := true
		for key_queue_full {
			time.Sleep(time.Millisecond * 50)
			key_queue_data, err := server.PeekMemory(KEYBOARD_QUEUED_KEYS_OFS, 1)
			if err != nil {
				time.Sleep(time.Second * 2)
				continue
			}
			keys_in_queue := key_queue_data[0]

			key_queue_full = keys_in_queue > 1
		}
	}

	return nil
}

func TypeChrCode(server server.U64Server, key byte) error {
	var c64_key []byte
	c64_key = append(c64_key, key)

	err := server.PokeMemory(KEYBOARD_BUFFER_OFS, c64_key)
	if err != nil { return err }

	err  = server.PokeMemory(KEYBOARD_QUEUED_KEYS_OFS, []byte{1})
	if err != nil { return err }

	return nil
}

func AppendChrCode(server server.U64Server, key byte) error {
	var c64_key []byte
	c64_key = append(c64_key, key)


	err := server.PokeMemory(KEYBOARD_BUFFER_OFS, c64_key)
	if err != nil { return err }

	err  = server.PokeMemory(KEYBOARD_QUEUED_KEYS_OFS, []byte{1})
	if err != nil { return err }

	return nil
}

func UnicodeToPet(char rune) byte {
	c64_key, has_char := priv_unicode_remap[char]

	if has_char {
		return c64_key
	} else {
		if unicode.IsLetter(char) {
			return byte(char)
		}
	}

	return priv_unicode_remap['?']
}


func TypeString(server server.U64Server, str string) error {
	var c64_keys []byte
	for _, char := range str {
		c64_key := UnicodeToPet(char)

		c64_keys = append(c64_keys, c64_key)
	}

	return TypeBytes(server, c64_keys)
}
