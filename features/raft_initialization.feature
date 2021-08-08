Feature: When the raft node initializes

  It should take care of the persistentState, should check the provided config, and after initialization, election should be in order if leader heart is not received

  Scenario: Raft node initializes for the first time on a system
    Given raft node initializes
    When Persistent state is not available
    And no evidence of crash is found
    Then raft node should become a follower
    And start the leader heartbeat monitor
    And reset currentTerm
    And reset logs
    And reset voteFor
    And reset CommitLength

  Scenario: Raft node initializes from a crash and persistent state is available
    Given raft node initializes
    When Persistent state is not available
    And evidence of crash is found
    Then raft node should become a follower
    And load up currentTerm from persistent storage 
    And should load up logs from persistent storage 
    And should load up VotedFor from persistent storage 
    And should load up CommitLength from persistent storage

  Scenario: Raft node initializes from a crash and persistent state is not available
    Given raft node initializes
    When Persistent state is not available
    And evidence of crash is found
    Then raft node should become a follower
    And start the leader heartbeat monitor
    And reset currentTerm
    And reset logs
    And reset voteFor
    And reset CommitLength


