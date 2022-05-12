$(document).ready(function() {
    setTimeout(function() {
        $(".alert").alert('close');
    }, 2000);

    $('#record_sermon').on('click', function() {
        alert($(this).prop('checked'));
        $('#config').submit();
    });

    $('#record_band').on('click', function() {
        alert($(this).prop('checked'));
        $('#config').submit();
    });

    function strcmp(a, b) {
        // return (a<b?-1:(a>b?1:0));
        return (a === b);
    }

    // var lastCmd = "";
    // var lastArg = "";
    function refresh_trigger() {
        // var lc = "";
        // $.get("/lastcmd", function( my_var ) {
        //     lc = my_var.toString();
        //     $("#LastCmd").text(lc.toString());
        // });
        // var la = "";
        // $.get("/lastargs", function( my_var ) {
        //     la = my_var;
        //     $("#LastArgs").text(la.toString());
        // });

        $("#LastArgs").load("/lastargs");
        args = $("#LastArgs").text();
        // args = args.replace(/[\n\r]+/g, '').trim();

        $("#LastCmd").load("/lastcmd");
        cmd = $("#LastCmd").text();
        // cmd = cmd.replace(/[\n\r]+/g, '').trim();

        // console.log("CMD:" + cmd + "ARGS:" + args + ": lastCmd:" + lastCmd + ": lastArg:" + lastArg + ":");
        //
        // if (strcmp(args, lastArg) && strcmp(cmd, lastCmd)) {
        //     return;
        // }
        //
        // if (strcmp(args, "NA")) {
        //     return;
        // }
        //
        // lastCmd = cmd;
        // lastArg = args;
        //
        // console.log("LA:" + args + ": lastCmd:" + lastCmd + ": lastArg:" + lastArg + ":");

        if (strcmp(args, "ON")) {
            $("#ON").addClass("highlight");
            $("#OFF").removeClass("highlight");
            $('input[name="Trigger"]').val(cmd);
        } else if (strcmp(args, "OFF")) {
            $("#ON").removeClass("highlight");
            $("#OFF").addClass("highlight");
            $('input[name="Trigger"]').val(cmd);
        } else {
            $("#ON").removeClass("highlight");
            $("#OFF").removeClass("highlight");
            $('input[name="Trigger"]').val("");
        }

        // $("#Trigger").text(cmd);
        // document.getElementById("Trigger").value = cmd;
    }
    setInterval(function() { refresh_trigger(); }, 1000);
});
