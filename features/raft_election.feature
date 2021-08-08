Feature: Raft election

  Raft election starts when a follower node doesn't receive a leader's heartbeat for a specified random duration
  Elections should be concluded within the election timeout duration, else the elections are restarted.
  
  Scenario: Raft node should start an election when it doesn't hear from the leader
    Given Raftnode has initialized
    When raft node doesn't receive leader's heartbeat
    Then it should start an election
    And become a candidate
    And vote for itself
    And increment the currentTerm
    And request votes from peers
  
  Scenario: Raft node should restart an election when it overtimes
    Given Raftnode has initialized
    When raft node doesn't receive leader's heartbeat
    Then it should start an election
    And vote for itself
    And increment the currentTerm
    And request votes from peers
    When election cannot reach to a conclusion within the allocated time
    Then election should be stopped
    And it should start an election
    And vote for itself
    And increment the currentTerm
    And request votes from peers
  
  Scenario: Raft node starts an election and receives majority votes
    Given Raftnode has initialized
    When raft node doesn't receive leader's heartbeat
    Then it should start an election
    And vote for itself
    And increment the currentTerm
    And request votes from peers
    When majority of peers vote in favor
    Then it should become a leader
    And stop the election
    And replicate logs to all its peers

  Scenario: Raft node starts an election and receives a vote with higher term
    Given Raftnode has initialized
    When raft node doesn't receive leader's heartbeat
    Then it should start an election
    And vote for itself
    And increment the currentTerm
    And request votes from peers
    When receives a vote with higher term
    Then it should become a follower
    And its current term should be equal to the higher term received
    And replicate logs to all its peers
    And stop the election
    When raft node doesn't receive leader's heartbeat
    Then it should start an election

  Scenario: Raft node receives a vote request and candidate's logs and term are ok, it should vote yes
    Given there are three raft nodes initialized
    And one of them starts an election
    And one of the node receives the vote request
    Then it should check if candidate's logs are ok
    And it should check if candidate's term is ok
    When candidate's logs and term are ok
    Then it should vote yes

  Scenario: Raft node receives a vote request and candidate's logs and term are not ok, it should vote no
    Given there are three raft nodes initialized
    And one of them starts an election
    And one of the node receives the vote request
    Then it should check if candidate's logs are ok
    And it should check if candidate's term is ok
    When candidate's logs and term are not ok
    Then it should vote no

  Scenario: Raft node receives a vote request and candidate's logs and term are not ok, it should vote no
    Given there are three raft nodes initialized
    And one of them starts an election
    And one of the node receives the vote request
    Then it should check if candidate's logs are ok
    And it should check if candidate's term is ok
    When candidate's logs and term are not ok
    Then it should vote no

  Scenario: Raft node starts an election, and during the election it receives a leader's heartbeat
    Given there are five raft nodes initialized
    Given one of the node is a leader
    When one of the non leader node does not receive leader's heartbeat for a while
    Then it starts an election
    And receives the leader heartbeat during the election
    Then it should become a follower
    And make term equal to that of leader

  Scenario: Raft node starts an election, and during the election it receives another candidate's vote request
    Given there are five raft nodes initialized
    When atleast two nodes start the election close enough
    Then 
    When one of the candidate receives another candidate's vote request
    Then ...
