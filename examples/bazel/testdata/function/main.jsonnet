local f = import "./function.libsonnet";

{
    x: f.doStuff(
        {
            hello: "world"
        }
    )
}