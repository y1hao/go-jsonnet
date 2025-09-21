local two = import "./2.jsonnet";

{
    convert:: function(x) (if x == 1 then 'one' else two.convert(x) )
}