const btn = document.querySelector("#btn");
const btnText = document.querySelector("#btnText");


btn.onclick = () => {
    var uname = document.forms["usersform"]["username"].value;
    var task = document.forms["usersform"]["tasks"].value;

    if (uname == "")
{
alert("Please enter a username");
} else if(task == "") {
    alert("action not defined")
} else {
    document.getElementById('usersform').submit();

    btnText.innerHTML = "process";
    btn.classList.add("active");
}
};