$(function(){
    var bgPath = ASSETS_X_URL+'/images/background/', imgs = ["1.jpg","2.jpg","3.jpg","4.jpg","5.jpg","6.jpg","7.jpg","8.jpg","9.jpg","10.jpg","11.jpg","12.jpg"];
    setInterval(function(){
    var index = Math.floor(Math.random() * imgs.length);
    $('#sectionSignIn').css("background-image","url('"+bgPath+imgs[index]+"')");
    },5000)
})