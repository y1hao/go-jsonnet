local four = import "./4.jsonnet";

{
    convert:: function(x) (if x == 3 then 'three' else four.convert(x) )
}