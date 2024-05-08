# Problem Statement
Design a messenger platform that can serve for around more than 100M users. The platform supports web and mobile apps(android, ios).
# Gathering Requirements
- Chatting:
  + Support one-on-one conversation.
  + Support group chat that has more than 1000 users per group.
  + Support those who already read the last message.
  + Support message sent time and message read time per user.
- Sending file, and attachment in a conversation:
  + Maximum size of a file, attachment is 5MB.
  + Maximum per sending is 5 files.
  + Just support image types(jpg, jpeg, png, SVG, ...) or document types(doc, docx, csv, xlsx).
  + Support emoji reaction for each message.
- Story video:
  + User can post brief videos.
  + User can post 10 brief videos per day.
  + Video duration must not exceed 60 seconds.
  + Support emoji reaction for each video.
# Capacity Planning
- Number of users on the application every month: 1M
- Number of active users per second at peak traffic: 500K
- Number of messages per second at peak traffic: 10M
- Estimate Database store size: 1M users * 5KB per user
# High-level Design
## Chatting
![chatting_high_level_design drawio](https://github.com/duyledat197/messenger/assets/38448882/c592f2b1-18d3-4dc6-8e95-8af82e02eeec)

## Emoji
![messenger_high_level_design](https://github.com/duyledat197/messenger/assets/38448882/a6af24ec-0cf7-44e4-8caa-6b58d5468066)



# Protocols
- Pros/Cons Between databases
- Decided or Plan
# Database
- Pros/Cons Between databases
- Decided or Plan

# API Design

# Monitoring
1. Availability SLA: p99.999
2. Latency SLA: p99.99 of 5 milliseconds
