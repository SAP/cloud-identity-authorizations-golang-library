// TEST FOR DCL FEATURES
schema {
/*
    stringval: String,
    $same:   {
        x: String
    },
    $Struct: {
        "Quoted-sub-name": String
    },
    "\"quoted\"": String,
    "\"quoted2\"": {
        findme : String
    },
    default: String,
    SalesOrders: String, 
    SalesOrderItems: String, 
    CountryCode: String,
    SalesId: Number,
    ActionPossible: Boolean,
    Company: {
        id: Number,
        name: String
    },
*/
    LEDGER: String,
    COMPANYCODE : String,
    COSTCENTER : String,
    SUPPLIER_AUTHORIZATIONGROUP : String,
    GLACCOUNTINCOMPANYCODE_AUTHORIZATIONGROUP : String,
    CUSTOMER_AUTHORIZATIONGROUP : String,
    ACCOUNTINGDOCUMENTTYPE_AUTHORIZATIONGROUP : String,
    PROFITCENTER : String,
    PROFITCTRRESPONSIBLEUSER: String,
    CONTROLLINGAREA : String,
    BUSINESSAREA : String,
    SEGMENT : String,
    void : Boolean,
    PostingDate : Number,
    FINANCIALACCOUNTTYPE: String,
    COSTCTRRESPONSIBLEUSER: String,
    PLANT: String,
    SALESORGANIZATION: String,
    DSITRIBUTIONCHANNEL: String,
    SalesDocumentType: String,
    ServiceDocumentType: String,
    ASSETCLASS: String,
    ORDERTYPE: String,
    VALUATIONAREA: String,
    $user: {
        id: String,
        user_uuid: String,
        "ledger": String,
         company: String,
         profitcenter: String,
         costcenter: String,
         controllingarea: String,
         authgrp: String,
         authgrp2: String,
         authgrp3: String,
         authgrp4: String,
         businessarea: String,
         segment: String,
         FAUTHCNTXT1 : Boolean,
         FAUTHCNTXT2 : Boolean, 
         salesorg : String,
         distributionchannel : String,
         SalesDocumentType: String,
         ServiceDocumentType: String,  
         ASSETCLASS: String,
         VALUATIONAREA: String,
         plant: String,
        "PersNumber": Number
    }
}