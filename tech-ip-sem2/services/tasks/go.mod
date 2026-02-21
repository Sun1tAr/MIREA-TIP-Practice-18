module github.com/sun1tar/MIREA-TIP-Practice-18/tech-ip-sem2/tasks

go 1.22

require (
    github.com/sun1tar/MIREA-TIP-Practice-18/tech-ip-sem2/proto v0.0.0
    github.com/sun1tar/MIREA-TIP-Practice-18/tech-ip-sem2/shared v0.0.0
    google.golang.org/grpc v1.64.0
)

replace (
    github.com/sun1tar/MIREA-TIP-Practice-18/tech-ip-sem2/proto => ../../proto
    github.com/sun1tar/MIREA-TIP-Practice-18/tech-ip-sem2/shared => ../../shared
)