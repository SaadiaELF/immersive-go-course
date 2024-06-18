Test our implementation and observe both of our consumers running jobs scheduled by your producer. 
- What happens if we only create one partition in our topic? 
=> All messages are consumed by only one consumer
- What happens if we create three?
=> Messages from partition 0 and 1 were consumed by consumer 1 and message from partition 2 were consumed by consumer 2

