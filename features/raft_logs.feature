//TODO:: convert ginko bdd to gherkin

	Context("Broadcasting Messages", func() {
		When("A Leader receives a broadcast request", func() {

			X("Should Append the entry to its log"})

			X("Should update the Acknowledged log length for itself"})

			X("Should Replicate the log to its followers after appending to its own logs"})
		})

		When("A Raft node is a leader", func() {
			X("Should send heart beats to its followers"})
		})

		When("A Follower receives a broadcast request", func() {
			X("Should forward the request to the leader node"})
		})
	})

	Context("Log Replication", func() {
		When("Leader is replicating log to followers", func() {
			Specify(`The log request should have leaderId, 
			currentTerm, 
			index of the last sent log entry,
			term of the previous log entry,
			commitLength and
			suffix of log entries`, func() {

			})
		})

		When("Follower receives a log replication request", func() {
			When("the received log is ok", func() {

				X("Should update its currentTerm to to new term"})

				X("Should continue being the follower and cancel the election it started (if any)"})

				X("Should update its acknowledged length of the log"})

				X("Should append the entry to its log"})

				X("Should send the log response back to the leader with the acknowledged new length"})

			})

			When("The received log or term is not ok", func() {

				XIt(`Should send its currentTerm
			acknowledged length as 0
			acknowledgment as false to the leader`, func() {

				})
			})

		})
	})

	Context("A node updates its log based on the log received from leader", func() {
		//TODO
	})

	Context("Leader Receives log acknowledgments", func() {
		When("The term in the log acknowledgment is more than leader's currentTerm", func() {})
		When("The term in the log acknowledgment is ok and log replication has been acknowledged by the follower", func() {

			When("log replication is acknowledged by the follower", func() {

				XIt(`Should update its acknowledged and 
					sentLength for the follower`, func() {

				})

				X("Should commit the log entries to persistent storage"})
			})

			When("log replication is not acknowledged by the follower", func() {

				X("Should update the sent length for the follower to one less than the previously attempted sent length value of the log"})
				X("Should send a log replication request to the follower with log length = last request log length - 1"})
			})

		})

	})

	Context("Commiting log entries to persistent storage", func() {
		//TODO
		When("The commit is successful it should send the log message to the client of Sif", func() {

			X("Don't know yet what to do""Not Yet Implemented")
			})
		})
	})
