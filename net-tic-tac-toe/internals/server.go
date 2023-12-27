

func startClient() {
	var err error

	rdr := bufio.NewReader(os.Stdin)
	var msg string

	for {
		fmt.Printf("Huh?: ")
		msg, err = rdr.ReadString('\n')
		if err != nil {
			panic(err)
		}
		msg = strings.TrimSpace(msg)

		if len(msg) > 0 {
			conn, err := net.Dial("tcp", ":8080")
			if err != nil {
				panic(err)
			}

			fmt.Fprint(conn, msg)
		}
	}
}
