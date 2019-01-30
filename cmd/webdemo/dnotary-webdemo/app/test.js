var helpService = "/api"
var mainChainEndpoint = "http://127.0.0.1:19888"
var sideChainEndpoint = "http://127.0.0.1:17888"
var pagetimer;
var key = localStorage["mykey"]; //秘钥
var nodes_Ip; //公证人列表
var myaccount = localStorage["myaccount"]; //当前账户地址
var mainChainContract; //主链合约地址
var sideChainContract;//侧链合约地址
var currentLockinSecretHash=localStorage["currentLockinSecretHash"];// 当前锁定hash
var currentLockinSecret=localStorage["currentLockinSecret"];//当前密码
var currentLockoutSecret=localStorage["currentLockoutSecret"]
var currentLockoutSecretHash=localStorage["currentLockoutSecretHash"]
var notaryPrivateKeyId; //公证人操作合约使用的私钥编号
var currentMainChainBlockNumber;
var currentSideChainBlockNuber
// localStorage.removeItem("myaccount")
$(function () {
    $(".btcpanel").hide();
    $('#tab_content').hide();
    $("#cointype").hide()
    nodes_Ip = [{
        name: 'Notary0',
        value: '127.0.0.1:8030'
    }, {
        name: 'Notary1',
        value: '127.0.0.1:8031'
    }, {
        name: 'Notary2',
        value: '127.0.0.1:8032'
    }, {
        name: 'Notary3',
        value: '127.0.0.1:8033'
    }, {
        name: 'Notary4',
        value: '127.0.0.1:8034'
    }, {
        name: 'Notary5',
        value: '127.0.0.1:8035'
    }, {
        name: 'Notary6',
        value: '127.0.0.1:8036'
    },];

    $('#selNode').html();
    var html = '';
    for (var i = 0; i < nodes_Ip.length; i++) {
        html += '<option value="http://' + nodes_Ip[i].value + '">' + nodes_Ip[i].name + '</option>';
    }
    $('#selNode').html(html);

    $("#selCoin").val("ETH");

})

function createKey(obj) {
    $(obj).attr("disabled", "disabled");
    $("#btnCreateKey").attr("disabled", "disabled");

    $("#privateKey").attr("readonly", "readonly");
    key = new Bitcoin.ECKey(false);
    key.setCompressed(false);
    var pubkey = key.getPubKeyHex();
    var privatekey = key.getBitcoinHexFormat();
    $("#privateKey").val(privatekey);
    $("#pubKey").val(pubkey); //todo 这里需要一个根据pubkey算地址的接口,
    if (privatekey) {
        $("#priv_warn").show();
        $('#tab_content').show();
    }
    $.ajax({
            url: helpService + "/pubkey2address/" + $("#pubKey").val(),
            type: "get",
            contentType: 'application/json',
            success: function (data) {
                console.log(data);
                if (data.Error != "") {
                    console.log("error", data.Error);
                    showTip("Error:" + JSON.stringify(data) + '<br/><br/>Please Retry!');
                    return
                }
                $("#address").html(data.Message)
                myaccount = data.Message
                localStorage["myaccount"] = myaccount
                localStorage["mykey"] = key.getBitcoinWalletImportFormat()
                queryStatus()
            },
            error: function (e) {
                console.log("error", e);
                showTip("Error:" + JSON.stringify(e) + '<br/><br/>Please Retry!');
            }
        },
    )
}

function transfer10Ether(obj) {
    if (!myaccount) {
        alert("generate key first")
        return
    }
    showMaskLayer("transfer test ether and smt to this address,please wait...");
    $.ajax(
        {
            url: $("#selNode").val() + "/api/1/debug/transfer-to-account/" + myaccount,
            type: "get",
            success: function (data) {
                hideMaskLayer()
                queryStatus()
                alert("your account already have test token.")
            },
            error: function (e) {
                hideMaskLayer();
                console.log("error", e);
                showTip("Error:" + JSON.stringify(e) + '<br/><br/>Please Retry!');
            }
        })
}

function showMaskLayer(val) {
    $("#maskLayerModal .masktext").html(val);
    $("#maskLayerModal").modal({backdrop: 'static', keyboard: false});
    var tempIndex = 0;
    $("#pagetimer").text(tempIndex++);
    pagetimer = window.setInterval(function () {
        $("#pagetimer").text(tempIndex++);
    }, 1000);
}

function updateMaskLayer(val) {
    $("#maskLayerModal .masktext").html(val);
}

function hideMaskLayer() {
    window.clearInterval(pagetimer);
    $("#maskLayerModal").modal('hide');
}

function changecoin(obj) {
    if ($(obj).val() == "BTC") {
        $(".ethpanel").hide();
        $(".btcpanel").show();
        $('.coinspan').text('testnet btc');
        displayAddress($('#btcDCRMAddress')[0]);
    } else {
        $(".ethpanel").show();
        $(".btcpanel").hide();
        $('.coinspan').text('rinkeby ether');
        displayAddress($('#ethDCRMAddress')[0]);
    }
}
//-------------------------- tx prepare lockin
function prePareLockin(obj) {
    currentLockinSecret = ""
    currentLockinSecretHash = ""
    var myBalance = getMainChainWeb3().toWei($("#MainChainBalance")[0].innerText, "milli")
    var amount = Math.floor($("#prepareLockInAmount").val() * myBalance)
    if (amount <= 0) {
        alert("no enough ether to transfer")
        return
    }
    showMaskLayer("prepare transfer ether to spectrum, amount=" + amount)
    $("#signTransaction").text('');
    $.ajax({
        url: helpService + "/generateSecret",
        type: "get",
        success: function (r) {
            if (r.Error) {
                hideMaskLayer()
                console.log("error", r.Error);
                showTip("generateSecret Error:" + r.Error + '<br/><br/>Please Retry!');
                return
            }
            r = r.Message
            currentLockinSecret = r.Secret
            currentLockinSecretHash = r.SecretHash
            localStorage["currentLockinSecret"]=currentLockinSecret
            localStorage["currentLockinSecretHash"]=currentLockinSecretHash
            //进行下一步操作,构造PrePareLockin调用,更新块数,然后继续
            queryStatusHelper(function (r) {
                if (r.Error) {
                    hideMaskLayer()
                    console.log("error", r.Error);
                    showTip("query status  Error:" + r.Error + '<br/><br/>Please Retry!');
                    return
                }
                r = r.Message
                //1000块的过期时间
                doPrePareLockin(amount, currentLockinSecretHash, r.MainChainBlockNumber + 1000)
            }, function () {
                hideMaskLayer()
            })
        },
        error: function (e) {
            hideMaskLayer();
            console.log("error", e);
            showTip("generate Secret Error, Please Retry!");
        }
    })
}

//secret,secrethash已经生成好了,直接用吧.
function doPrePareLockin(amount, secrethash, expiration) {
    var req = {}
    req.From = myaccount
    req.ContractAddress = mainChainContract
    req.Method = "mprepareLockin"
    req.Arg = {
        SecretHash: secrethash,
        Expiration: expiration,
        Value: amount,
    }
    $.ajax(
        {
            url: helpService + "/generateTx",
            type: "post",
            dataType: "json",
            data: JSON.stringify(req),
            contentType: 'application/json',
            success: function (r) {
                if (r.Error) {
                    hideMaskLayer();
                    console.log("generateTx err ", r.Error);
                    showTip("generateTx Error:" + r.Error + '<br/><br/>Please Retry!');
                    return
                }
                r = r.Message
                console.log("message to sign:" + r.TxHash)

                $("#signTransaction").text(formatJson(JSON.stringify(r.Tx),
                    {
                        newlineAfterColonIfBeforeBraceOrBracket: true,
                        spaceAfterColon: true,
                    }
                ))
                doSendTx(r, "main",function(){
                    // alert("notify notary to assign ethereum token for you")
                    setTimeout(notifyNotaryPreareLockin,5000) //延时五秒执行,让公证人节点知道交易
                })
            },
            error: function (e) {
                hideMaskLayer();
                console.log("error", e);
                showTip("mprepareLockin Error:" + JSON.stringify(e) + '<br/><br/>Please Retry!');
            }

        }
    )
}
//-------------------------- tx prepare lockin
function  stripzero(s){
    var ss=""
    var i=0;
    for (;i<s.length;i++){
        if (s[i]!='0'){
            break
        }
    }
    if (i>0){
        console.log("has zero ")
        printCallStack()
    }
    return s.slice(i)
}
function printCallStack() {
    var i = 0;
    var fun = arguments.callee;
    do {
        fun = fun.arguments.callee.caller;
        console.log(++i + ': ' + fun);
    } while (fun);
}
/*
r: send tx req
chain: main 或者side
cb: tx成功以后执行的回调函数
 */
function doSendTx(r, chain,cb) {
    var req = {}
    req.Chain = chain
    req.Tx = r.Tx
    req.TxHash = r.TxHash
    req.Signer = myaccount
    var hasharray = Crypto.util.hexToBytes(r.TxHash);
    var signarray = key.sign(hasharray);
    var obj = key.parseSigHex(signarray);
    req.Tx.r = "0x" + stripzero(obj.r)
    req.Tx.s = "0x" + stripzero(obj.s)
    req.Tx.v = "0x0"
    updateMaskLayer("send tx to " + chain + "...")
    $.ajax(
        {
            url: helpService + "/sendTx",
            type: "post",
            data: JSON.stringify(req),
            contentType: "application/json",
            success: function (r) {
                if (r.Error) {
                    hideMaskLayer();
                    console.log("sendTx err ", r.Error);
                    showTip("generateTx Error:" + r.Error + '<br/><br/>Please Retry!');
                    return
                }
                //广播结果有了
                r = r.Message
                console.log("tx receipt " + JSON.stringify(r))
                queryStatus()
                if(cb){
                    cb()
                }
            },
            error: function (e) {
                hideMaskLayer();
                console.log("error", e);
                showTip("sendTx Error:" + JSON.stringify(e) + '<br/><br/>Please Retry!');
            }
        }
    )
}

//通知公证人侧链PrepareLockin,安全的实现,需要等待块数到达以后,目前不用.
function notifyNotaryPreareLockin(obj) {
    var notary = $("#selNode").val()
    var req = {}
    req.SCToken = sideChainContract
    req.UserAddress = myaccount
    req.SecretHash = currentLockinSecretHash
    if (!req.SecretHash){
        alert("please prepare lock in first")
        return
    }

    updateMaskLayer("notify notary "+notary+" assign Ethereum Token for me on Spectrum...")
    //使用helpService服务构造发给公证人的PrePareLockin请求,在js断构并计算hash会出问题
    $.ajax({
        url: helpService + "/scPrepareLockin",
        type: "post",
        dataType: "json",
        data: JSON.stringify(req),
        contentType: "application/json",
        success: function (r) {
            if(r.Error){
                hideMaskLayer();
                console.log("scPrepareLockin err ", r.Error);
                showTip("scPrepareLockin Error:" + r.Error + '<br/><br/>Please Retry!');
                return
            }
            r = r.Message
            console.log("scPrepareLockin help service:  " + JSON.stringify(r))
            doNotifyNotaryPrepareLockin(r)
        },
        err: function (e) {
            hideMaskLayer()
            console.log("query status err ", e.responseText)
            showTip("[Error] " + JSON.parse(e.responseText).error + '<br/><br/>Please Retry!');
        }
    })
}
function doNotifyNotaryPrepareLockin(rFromHelpService) {
    var notary = $("#selNode").val()
    var r=rFromHelpService
    //签名发给公证人的请求,所有签名必须发生在网页端
    var hasharray = Crypto.util.hexToBytes(r.TxHash);
    var signarray = key.sign(hasharray);
    var obj = key.parseSigHex(signarray);
    var rsv=obj.r+obj.s+"00"
    rsv=Crypto.util.hexToBytes(rsv)
    var rsvBase64=Crypto.util.bytesToBase64(rsv)
    r.Req.signature=rsvBase64
    $.ajax({
        url:notary+"/api/1/user/scpreparelockin/"+sideChainContract,
        type:"post",
        dataType:"json",
        data:JSON.stringify(r.Req),
        contentType:"application/json",
        success:function(r){
            if(r.error_msg!="success"){
                hideMaskLayer();
                console.log("scPrepareLockin notary err ", JSON.stringify(r));
                showTip("scPrepareLockin notary Error:" +JSON.stringify(r) + '<br/><br/>Please Retry!');
                return
            }
            r=r.Data
            console.log("scPrepareLockin help service:  " + JSON.stringify(r))
            queryStatus()
            //lockin for side chain
            sideChainLockin()
        },
        err:function(e){
            hideMaskLayer()
            console.log("query status err ", e.responseText)
            showTip("[Error] " + JSON.parse(e.responseText).error + '<br/><br/>Please Retry!');
        }
    })
}

function queryStatusHelper(cb, cberror) {
    var req = {
        MainChainContract: mainChainContract,
        SideChainContract: sideChainContract,
        Account: myaccount,
        LockSecretHash: currentLockinSecretHash,
    }
    $.ajax({
        url: helpService + "/querystatus",
        type: 'post',
        dataType: 'json',
        data: JSON.stringify(req),
        contentType: 'application/json',
        success: function (data) {
            if (cb) {
                cb(data)
            }

        },
        err: function (e) {
            console.log("query status err ", e.responseText)
            showTip("[Error] " + JSON.parse(e.responseText).error + '<br/><br/>Please Retry!');
            if (cberror) {
                cberror()
            }

        }
    })
}

function queryStatus() {
    queryStatusHelper(function (r) {
        if (r.Error) {
            console.log("query status err ", r.Error)
            showTip("[Error] " + r.Error + '<br/><br/>Please Retry!');
            return
        } else {
            r = r.Message
        }
        currentMainChainBlockNumber = r.MainChainBlockNumber
        currentSideChainBlockNuber = r.SideChainBlockNumber
        $("#MainChainBlockNumber").html(r.MainChainBlockNumber)
        $("#SideChainBlockNumber").html(r.SideChainBlockNumber)
        $("#mainChainContractBalance").html(getMainChainWeb3().fromWei(r.MainChainContractBalance, "ether"))
        $("#sideChainContractBalance").html(getMainChainWeb3().fromWei(r.SideChainContractBalance, "ether"))
        $("#MainChainBalance").html(getMainChainWeb3().fromWei(r.MainChainBalance, "milli"))
        $("#SideChainBalance").html(getMainChainWeb3().fromWei(r.SideChainBalance, "milli"))
        $("#SideChainTokenBalance").html(getMainChainWeb3().fromWei(r.SideChainTokenBalance, "milli"))
        if (r.MainChainBalance>0) {
            $("#btnTransferEther").attr("disabled", "disabled");
        }
        //判断一下是lockin还是lockout
        if (localStorage["currentLockinSecret"]){
            $("#btnPrepareLockin").attr("disabled",true)
            $("#btnMainChainCancelLockin").attr("disabled",false)
        } else{
            $("#btnPrepareLockin").attr("disabled",false)
            $("#btnMainChainCancelLockin").attr("disabled",true)
        }
        if (localStorage["currentLockoutSecret"]){
            $("#btnPrepareLockout").attr("disabled",true)
            $("#btnSideChainCancelLockout").attr("disabled",false)
        } else{
            $("#btnPrepareLockout").attr("disabled",false)
            $("#btnSideChainCancelLockout").attr("disabled",true)
        }
    })
}

//初始化状态
$(function () {
    var notary = $("#selNode").val()
    $.ajax({
        url: notary + "/api/1/user/sctokens",
        type: "get",
        contentType: "application/json",
        success: function (r) {
            if (r.error_msg != "success") {
                console.log("contract list  err ", JSON.stringify(r))
                showTip("[Error] " + JSON.stringify(r).error + '<br/><br/>Please Retry!');
                return
            }
            r = r.Data[0] //假定只有一组合约,后续需要完善
            mainChainContract = r.mc_locked_contract_address
            sideChainContract = r.sc_token
            notaryPrivateKeyId = r.sc_token_owner_key
            $("#mainChainContract").html(mainChainContract)
            $("#sideChainContract").html(sideChainContract)
            //合约地址有了,更新状态吧.
            queryStatus()
        },
        err: function (e) {
            console.log("query status err ", e.responseText)
            showTip("[Error] " + JSON.parse(e.responseText).error + '<br/><br/>Please Retry!');
        }
    })
    //如果以前有key,直接拿过来用
    if (myaccount && myaccount.length > 0) {
        $("#btnCreateKey").attr("disabled", "disabled");

        $("#privateKey").attr("readonly", "readonly");
        key = new Bitcoin.ECKey(localStorage["mykey"])
        $("#privateKey").val(key.getBitcoinHexFormat())
        $("#address").html(myaccount)
        $('#tab_content').show();
    }

})
var mainWeb3;
var sideWeb3;

function getMainChainWeb3() {
    var Web3 = require('web3');

    if (mainWeb3) {
        return mainWeb3
    } else {
        mainWeb3 = new Web3(new Web3.providers.HttpProvider(mainChainEndpoint));
    }
    return mainWeb3

}

function getMainSideChainWeb3() {
    var Web3 = require('web3');

    if (sideWeb3) {
        return sideWeb3;
    } else {
        sideWeb3 = new Web3(new Web3.providers.HttpProvider(sideChainEndpoint));
    }
    return sideWeb3
}

//baidu
var _hmt = _hmt || [];
(function () {
    var hm = document.createElement("script");
    hm.src = "https://hm.baidu.com/hm.js?c1d8355244d33be9abb1c0c996889a9a";
    var s = document.getElementsByTagName("script")[0];
    s.parentNode.insertBefore(hm, s);
})();

var formatJson = function (json, options) {
    var reg = null,
        formatted = '',
        pad = 0,
        PADDING = '    '; // one can also use '\t' or a different number of spaces
    // optional settings
    options = options || {};
    // remove newline where '{' or '[' follows ':'
    options.newlineAfterColonIfBeforeBraceOrBracket = (options.newlineAfterColonIfBeforeBraceOrBracket === true) ? true : false;
    // use a space after a colon
    options.spaceAfterColon = (options.spaceAfterColon === false) ? false : true;

    // begin formatting...

    // make sure we start with the JSON as a string
    if (typeof json !== 'string') {
        json = JSON.stringify(json);
    }
    // parse and stringify in order to remove extra whitespace
    json = JSON.parse(json);
    json = JSON.stringify(json);

    // add newline before and after curly braces
    reg = /([\{\}])/g;
    json = json.replace(reg, '\r\n$1\r\n');

    // add newline before and after square brackets
    reg = /([\[\]])/g;
    json = json.replace(reg, '\r\n$1\r\n');

    // add newline after comma
    reg = /(\,)/g;
    json = json.replace(reg, '$1\r\n');

    // remove multiple newlines
    reg = /(\r\n\r\n)/g;
    json = json.replace(reg, '\r\n');

    // remove newlines before commas
    reg = /\r\n\,/g;
    json = json.replace(reg, ',');

    // optional formatting...
    if (!options.newlineAfterColonIfBeforeBraceOrBracket) {
        reg = /\:\r\n\{/g;
        json = json.replace(reg, ':{');
        reg = /\:\r\n\[/g;
        json = json.replace(reg, ':[');
    }
    if (options.spaceAfterColon) {
        reg = /\:/g;
        json = json.replace(reg, ': ');
    }

    $.each(json.split('\r\n'), function (index, node) {
        var i = 0,
            indent = 0,
            padding = '';

        if (node.match(/\{$/) || node.match(/\[$/)) {
            indent = 1;
        } else if (node.match(/\}/) || node.match(/\]/)) {
            if (pad !== 0) {
                pad -= 1;
            }
        } else {
            indent = 0;
        }

        for (i = 0; i < pad; i++) {
            padding += PADDING;
        }

        formatted += padding + node + '\r\n';
        pad += indent;
    });

    return formatted;
};
function clearData(){
    localStorage.removeItem("myaccount")
    localStorage.removeItem("mykey")
    clearLockinSecret()
    clearLockoutSecret()
}


//-------------------------- tx  lockin
function sideChainLockin(obj) {
    if(!currentLockinSecret || !currentLockinSecretHash) {
        alert("must prepare lockin and notify notary first")
        return
    }
    updateMaskLayer("get Ethereum Token on Spectrum,please wait ...")
    $("#signTransaction").text('');
    doSideChainLockin()
}

//secret,secrethash已经生成好了,直接用吧.
function doSideChainLockin( ) {
    var req = {}
    req.From = myaccount
    req.ContractAddress = sideChainContract
    req.Method = "slockin"
    req.Arg = {
        Account: myaccount,
        Secret:currentLockinSecret,
    }
    $.ajax(
        {
            url: helpService + "/generateTx",
            type: "post",
            dataType: "json",
            data: JSON.stringify(req),
            contentType: 'application/json',
            success: function (r) {
                if (r.Error) {
                    hideMaskLayer();
                    console.log("generateTx err ", r.Error);
                    showTip("generateTx Error:" + r.Error + '<br/><br/>Please Retry!');
                    return
                }
                r = r.Message
                console.log("message to sign:" + r.TxHash)

                $("#signTransaction").text(formatJson(JSON.stringify(r.Tx),
                    {
                        newlineAfterColonIfBeforeBraceOrBracket: true,
                        spaceAfterColon: true,
                    }
                ))
                //由helpService在侧连上执行Tx
                doSendTx(r, "side",function(){
                    alert("your eth have been moved to spectrum as EthereumToken")
                    clearLockinSecret()
                    hideMaskLayer()
                })
            },
            error: function (e) {
                hideMaskLayer();
                console.log("error", e);
                showTip("doSideChainLockin Error:" + JSON.stringify(e) + '<br/><br/>Please Retry!');
            }

        }
    )
}

//-------------------------- tx   lockin

//-------------------------- tx  cancel lockin
function mainChainCancelLockin(obj) {
    if(!currentLockinSecret || !currentLockinSecretHash) {
        alert("must prepare lockin and notify notary first")
        return
    }
    showMaskLayer("cancel Ethereum lockin,please wait ...")
    $("#signTransaction").text('');
    doMainChainCancelLockin()
}

//secret,secrethash已经生成好了,直接用吧.
function doMainChainCancelLockin( ) {
    var req = {}
    req.From = myaccount
    req.ContractAddress = mainChainContract
    req.Method = "mcancelLockin"
    req.Arg = {
        Account: myaccount,
    }
    $.ajax(
        {
            url: helpService + "/generateTx",
            type: "post",
            dataType: "json",
            data: JSON.stringify(req),
            contentType: 'application/json',
            success: function (r) {
                if (r.Error) {
                    hideMaskLayer();
                    console.log("generateTx err ", r.Error);
                    showTip("generateTx Error:" + r.Error + '<br/><br/>Please Retry!');
                    return
                }
                r = r.Message
                console.log("message to sign:" + r.TxHash)

                $("#signTransaction").text(formatJson(JSON.stringify(r.Tx),
                    {
                        newlineAfterColonIfBeforeBraceOrBracket: true,
                        spaceAfterColon: true,
                    }
                ))
                //由helpService在侧连上执行Tx
                doSendTx(r, "main",function(){
                    alert("your eth have been returned to your account")
                    clearLockinSecret()
                })
                //tx成功以后回调query status,然后,
                // queryStatus()
            },
            error: function (e) {
                hideMaskLayer();
                console.log("error", e);
                showTip("doMainChainCancelLockin Error:" + JSON.stringify(e) + '<br/><br/>Please Retry!');
            }

        }
    )
}

//-------------------------- tx cancel lockin

//-------------------------- tx  cancel lockout
function sideChainCancelLockout(obj) {
    if(!currentLockinSecret || !currentLockinSecretHash) {
        alert("must prepare lockout and notify notary first")
        return
    }
    showMaskLayer("cancel spectrum lockout,please wait ...")
    $("#signTransaction").text('');
    doSideChainCancelLockout()
}

//secret,secrethash已经生成好了,直接用吧.
function doSideChainCancelLockout( ) {
    var req = {}
    req.From = myaccount
    req.ContractAddress = sideChainContract
    req.Method = "scancelLockOut"
    req.Arg = {
        Account: myaccount,
    }
    $.ajax(
        {
            url: helpService + "/generateTx",
            type: "post",
            dataType: "json",
            data: JSON.stringify(req),
            contentType: 'application/json',
            success: function (r) {
                if (r.Error) {
                    hideMaskLayer();
                    console.log("generateTx err ", r.Error);
                    showTip("generateTx Error:" + r.Error + '<br/><br/>Please Retry!');
                    return
                }
                r = r.Message
                console.log("message to sign:" + r.TxHash)

                $("#signTransaction").text(formatJson(JSON.stringify(r.Tx),
                    {
                        newlineAfterColonIfBeforeBraceOrBracket: true,
                        spaceAfterColon: true,
                    }
                ))
                //由helpService在侧连上执行Tx
                doSendTx(r, "side",function(){
                    alert("your EtherumToken have been returned to your account")
                    clearLockinSecret()
                })
                //tx成功以后回调query status,然后,
                // queryStatus()
            },
            error: function (e) {
                hideMaskLayer();
                console.log("error", e);
                showTip("doSideChainCancelLockout Error:" + JSON.stringify(e) + '<br/><br/>Please Retry!');
            }

        }
    )
}

//-------------------------- tx cancel lockout

//-------------------------- tx prepare lockout
function prePareLockout(obj) {
    currentLockoutSecret = ""
    currentLockoutSecretHash = ""
    var myBalance = getMainChainWeb3().toWei($("#SideChainTokenBalance")[0].innerText, "milli")
    var amount = Math.floor($("#prepareLockoutAmount").val() * myBalance)
    if (amount <= 0) {
        alert("no enough EtherumEther to transfer")
        return
    }
    showMaskLayer("prepare transfer ether from   spectrum to ether, amount=" + amount)
    $("#signTransaction").text('');
    $.ajax({
        url: helpService + "/generateSecret",
        type: "get",
        success: function (r) {
            if (r.Error) {
                hideMaskLayer()
                console.log("error", r.Error);
                showTip("generateSecret Error:" + r.Error + '<br/><br/>Please Retry!');
                return
            }
            r = r.Message
            currentLockoutSecret = r.Secret
            currentLockoutSecretHash = r.SecretHash
            localStorage["currentLockoutSecret"]=currentLockoutSecret
            localStorage["currentLockoutSecretHash"]=currentLockoutSecretHash
            //进行下一步操作,构造PrePareLockout调用,更新块数,然后继续
            queryStatusHelper(function (r) {
                if (r.Error) {
                    hideMaskLayer()
                    console.log("error", r.Error);
                    showTip("query status  Error:" + r.Error + '<br/><br/>Please Retry!');
                    return
                }
                r = r.Message
                //1000块的过期时间
                doPrePareLockout(amount, currentLockoutSecretHash, r.SideChainBlockNumber + 1000)
            }, function () {
                hideMaskLayer()
            })
        },
        error: function (e) {
            hideMaskLayer();
            console.log("error", e);
            showTip("Error:" + JSON.stringify(e) + '<br/><br/>Please Retry!');
        }
    })
}

//secret,secrethash已经生成好了,直接用吧.
function doPrePareLockout(amount, secrethash, expiration) {
    var req = {}
    req.From = myaccount
    req.ContractAddress = mainChainContract
    req.Method = "sprepareLockout"
    req.Arg = {
        SecretHash: secrethash,
        Expiration: expiration,
        Value: amount,
    }
    $.ajax(
        {
            url: helpService + "/generateTx",
            type: "post",
            dataType: "json",
            data: JSON.stringify(req),
            contentType: 'application/json',
            success: function (r) {
                if (r.Error) {
                    hideMaskLayer();
                    console.log("generateTx err ", r.Error);
                    showTip("generateTx Error:" + r.Error + '<br/><br/>Please Retry!');
                    return
                }
                r = r.Message
                console.log("message to sign:" + r.TxHash)

                $("#signTransaction").text(formatJson(JSON.stringify(r.Tx),
                    {
                        newlineAfterColonIfBeforeBraceOrBracket: true,
                        spaceAfterColon: true,
                    }
                ))
                doSendTx(r, "side",function(){
                   // alert("notify the notary to  prepare your eth on Ethereum ")
                    setTimeout(notifyNotaryPreareLockout,5000)
                })
            },
            error: function (e) {
                hideMaskLayer();
                console.log("error", e);
                showTip("doPrePareLockout Error:" + JSON.stringify(e) + '<br/><br/>Please Retry!');
            }

        }
    )
}
//-------------------------- tx prepare lockin

//通知公证人侧链PrepareLockout,安全的实现,需要等待块数到达以后,目前不用.
function notifyNotaryPreareLockout(obj) {
    var notary = $("#selNode").val()
    var req = {}
    req.SCToken = sideChainContract
    req.UserAddress = myaccount
    req.SecretHash = currentLockoutSecretHash
    if (!req.SecretHash){
        alert("please prepare lock out first")
        return
    }

    updateMaskLayer("notify notary "+notary+" assign Eth for me on Etherum...")
    //使用helpService服务构造发给公证人的PrePareLockin请求,在js断构并计算hash会出问题
    $.ajax({
        url: helpService + "/mcPrepareLockout",
        type: "post",
        dataType: "json",
        data: JSON.stringify(req),
        contentType: "application/json",
        success: function (r) {
            if(r.Error){
                hideMaskLayer();
                console.log("mcPrepareLockout err ", r.Error);
                showTip("mcPrepareLockout Error:" + r.Error + '<br/><br/>Please Retry!');
                return
            }
            r = r.Message
            console.log("mcPrepareLockout help service:  " + JSON.stringify(r))
            doNotifyNotaryPrepareLockout(r)
        },
        err: function (e) {
            hideMaskLayer()
            console.log("query status err ", e.responseText)
            showTip("[Error] " + JSON.parse(e.responseText).error + '<br/><br/>Please Retry!');
        }
    })
}
function doNotifyNotaryPrepareLockout(rFromHelpService) {
    var notary = $("#selNode").val()
    var r=rFromHelpService
    //签名发给公证人的请求,所有签名必须发生在网页端
    var hasharray = Crypto.util.hexToBytes(r.TxHash);
    var signarray = key.sign(hasharray);
    var obj = key.parseSigHex(signarray);
    var rsv=obj.r+obj.s+"00"
    rsv=Crypto.util.hexToBytes(rsv)
    var rsvBase64=Crypto.util.bytesToBase64(rsv)
    r.Req.signature=rsvBase64
    $.ajax({
        url:notary+"/api/1/user/mcpreparelockout/"+sideChainContract,
        type:"post",
        dataType:"json",
        data:JSON.stringify(r.Req),
        contentType:"application/json",
        success:function(r){
            if(r.error_msg!="success"){
                hideMaskLayer();
                console.log("mcpreparelockout notary err ", JSON.stringify(r));
                showTip("mcpreparelockout notary Error:" +JSON.stringify(r) + '<br/><br/>Please Retry!');
                return
            }
            hideMaskLayer()
            r=r.Data
            console.log("mcpreparelockout help service:  " + JSON.stringify(r))
            queryStatus()
            sideChainLockout()

        },
        err:function(e){
            hideMaskLayer()
            console.log("query status err ", e.responseText)
            showTip("[Error] " + JSON.parse(e.responseText).error + '<br/><br/>Please Retry!');
        }
    })
}


//-------------------------- tx  lockin
function sideChainLockout(obj) {
    if(!currentLockinSecret || !currentLockinSecretHash) {
        alert("must prepare lockout and notify notary first")
        return
    }
    updateMaskLayer("get ETH Token on Ethereum,please wait ...")
    $("#signTransaction").text('');
    doSideChainLockout()
}

//secret,secrethash已经生成好了,直接用吧.
function doSideChainLockout( ) {
    var req = {}
    req.From = myaccount
    req.ContractAddress = sideChainContract
    req.Method = "mlockout"
    req.Arg = {
        Account: myaccount,
        Secret:currentLockoutSecret,
    }
    $.ajax(
        {
            url: helpService + "/generateTx",
            type: "post",
            dataType: "json",
            data: JSON.stringify(req),
            contentType: 'application/json',
            success: function (r) {
                if (r.Error) {
                    hideMaskLayer();
                    console.log("generateTx err ", r.Error);
                    showTip("generateTx Error:" + r.Error + '<br/><br/>Please Retry!');
                    return
                }
                r = r.Message
                console.log("message to sign:" + r.TxHash)

                $("#signTransaction").text(formatJson(JSON.stringify(r.Tx),
                    {
                        newlineAfterColonIfBeforeBraceOrBracket: true,
                        spaceAfterColon: true,
                    }
                ))
                //由helpService在侧连上执行Tx
                doSendTx(r, "main",function(){
                    alert("your EthereumToken have been moved to ethereum as eth")
                    clearLockoutSecret()
                    hideMaskLayer()
                })
            },
            error: function (e) {
                hideMaskLayer();
                console.log("error", e);
                showTip("doSideChainLockout Error:" + JSON.stringify(e) + '<br/><br/>Please Retry!');
            }

        }
    )
}

//-------------------------- tx   lockin
function clearLockinSecret(){
    localStorage.removeItem("currentLockinSecretHash")
    localStorage.removeItem("currentLockinSecret")
}
function clearLockoutSecret(){
    localStorage.removeItem("currentLockoutSecretHash")
    localStorage.removeItem("currentLockoutSecret")
}

window.onbeforeunload=function(e){
    var e = window.event||e;
    e.returnValue=("Leave this page？the data may error when reopen");
}