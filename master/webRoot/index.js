var host = "localhost:8080"

$(function () {
    // alert(1)
    refreshJobList()
})

function refreshJobList() {
    $(".job-list").html("")
    $.ajax({
        type:"POST",
        url:"/job/list",
        dataType:"json",
        success:function (data) {
            console.log(data)
            if(data.errNo != 0){
                alert("query error, error code:" + data.errNo)
                return
            }
            var jobs = data.data
            for(var i=0;i<jobs.length;i++){
                var job = jobs[i]
                buildJobList(job)
            }
        },
        error:function () {
            alert("query wrong")
        }
    })
}

function buildJobList(job) {
    var name = job.name
    var command = job.command
    var cronExpr = job.cronExpr

    var html = "<tr>" +
        "<td class=\"job-name\">"+name+"</td>" +
        "<td class=\"job-command\">"+command+"</td>" +
        "<td class=\"job-cronExpr\">"+cronExpr+"</td>" +
        "<td class=\"job-btns\">" +
        "<div class=\"btn-toolbar\">" +
        "<button class=\"btn btn-info edit-job\">编辑</button>" +
        "<button class=\"btn btn-danger delete-job\" onclick='javascript:deleteJob(\""+name+"\")'>删除</button>" +
        "<button class=\"btn btn-warning kill-job\" onclick='javascript:killJob(\""+name+"\")'>强杀</button>" +
        "</div>" +
        "</td>" +
        "</tr>"

    $(".job-list").append(html)
}

function deleteJob(jobName) {
    console.log("deleteJob: "+jobName)
    $.ajax({
        url:"/job/delete",
        type: "POST",
        data:'{"name":"'+jobName+'"}',
        dataType: "json",
        // contentType: "application/json; charset=utf-8",
        success:function (data) {
            console.log(data)
            if(data.errNo == 0){
                alert("删除成功")
                refreshJobList()
            }else{
                alert("删除失败,error: ",data.errNo)
            }
        },
        error:function () {
            alert("delete error")
        }
    })
}
function killJob(jobName) {
    console.log("killJob: "+jobName)
    $.ajax({
        url:"/job/kill",
        type: "POST",
        data:'{"name":"'+jobName+'"}',
        dataType: "json",
        // contentType: "application/json; charset=utf-8",
        success:function (data) {
            console.log(data)
            if(data.errNo == 0){
                alert("强杀成功")
                refreshJobList()
            }else{
                alert("强杀失败,error: ",data.errNo)
            }
        },
        error:function () {
            alert("delete error")
        }
    })
}