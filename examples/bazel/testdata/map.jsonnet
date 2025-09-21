local original = [
    {
        name: "a"
    },
    {
        name: "b"
    },
    {
        name: "c"
    }
];

std.map(function(x) (
    x {
        ["name of %s" % x.name]: x.name,
        greetings: "hello " + x.name
    }
), original)