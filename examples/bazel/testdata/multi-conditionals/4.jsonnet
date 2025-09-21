local five = import "./5.jsonnet";

{
    convert:: function(x) (if x == 4 then 'four' else five.convert(x) )
}