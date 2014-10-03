(* cgroup cgi thingie *)

#load "unix.cma"
#load "str.cma"

open Str

open Unix

let root = "/sys/fs/cgroup";;

let resources = Array.to_list (Sys.readdir root);;

let walk_dir dir =
    let rec walk acc = function
        | [] -> (acc)
        | dir :: tail ->
                let contents = Array.to_list (Sys.readdir dir) in
                let contents = List.rev_map (Filename.concat dir) contents in
                let ftype acc file =
                    match (stat file).st_kind with
                    | S_DIR -> file::acc
                    | _ ->  acc
                in
                let dirs =  List.fold_left ftype [] contents in
                walk (dirs @ acc) (dirs @ tail)
    in
    walk [dir ^ "/" ] [dir];;


let () =
    let print_resources res =
        let dir = root ^ "/" ^ res in
        let print_group path =
            let skip = String.length root + String.length res + 1 in
            let name = Str.string_after path skip in
            Printf.printf "%s:%s\n" res name
        in
        List.iter print_group (walk_dir dir)
    in
    List.iter print_resources resources;;
