
$("a[href^=http://report.intermesh.net/]").addClass("linkclass");

$(document).ready(function(){
    $(".linkclass").click(function(e){
        e.preventDefault();
        
        alert($(this).attr('href'));  // getting href value

        var objectData =
         {
             linktosend: $(this).attr('href')
                            
         };
         var objectDataString = JSON.stringify(objectData);
    // sending json data through ajax
    $.ajax({
            type: "POST",
            url: "/service",
            dataType: "json",
            data: {
                ajax_post_data: objectDataString
            },
            // success: function (data) {
            //    alert('Success');

            // },
            // error: function () {
            //  alert('Error');
            // }
        });

});
});
