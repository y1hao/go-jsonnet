local three = import "./3.jsonnet";

{
    convert:: function(x) (if x == 2 then 'two' else three.convert(x) )
}