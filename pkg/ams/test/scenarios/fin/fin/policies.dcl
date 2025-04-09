POLICY GLAccountLineItem {

  grant read on I_GLAccountLineItem where LEDGER = $user.ledger
    and
  (
       COMPANYCODE  = $user.company
    or
       COMPANYCODE  = '100'
  )
and
   GLACCOUNTINCOMPANYCODE_AUTHORIZATIONGROUP = $user.authgrp3
and
   SUPPLIER_AUTHORIZATIONGROUP  = $user.authgrp
and
  CUSTOMER_AUTHORIZATIONGROUP  = $user.authgrp2
and
  ACCOUNTINGDOCUMENTTYPE_AUTHORIZATIONGROUP = $user.authgrp4
and
  (
         PROFITCENTER  = $user.profitcenter
     or
        CONTROLLINGAREA  = $user.controllingarea
     or
         PROFITCTRRESPONSIBLEUSER = $user.id
    )  
and  PostingDate  BETWEEN 19990101 and 20212131
and
  (
      (
          BUSINESSAREA = $user.businessarea
        and
           SEGMENT  = $user.segment
        and
          (
              (
                   FINANCIALACCOUNTTYPE  = 'OPEN'
                and
                  (
                      void = $user.FAUTHCNTXT1
                    or
                      (
                          void = $user.FAUTHCNTXT2
                       and
                                                    (
                               CONTROLLINGAREA = $user.controllingarea
                            or
                              COSTCENTER = $user.costcenter
                            or
                              COSTCTRRESPONSIBLEUSER = $user.id
                          )
                        AND (
                              ( 
                                   SALESORGANIZATION = $user.salesorg
                               OR 
 
                                   DSITRIBUTIONCHANNEL  = $user.distributionchannel
                               )
                               AND SalesDocumentType  = $user.SalesDocumentType
                               AND ServiceDocumentType =  $user.ServiceDocumentType
                        )                         
                      )
                  )
              )
            or
              (
                  void = $user.FAUTHCNTXT1
                and
                  COMPANYCODE = $user.company
                and
                  BUSINESSAREA = $user.businessarea
              )
          )
      )
    or
      (
          ORDERTYPE =  'valid'
        and
          (
              (
                  void = $user.FAUTHCNTXT2 
                and
                  (
                      CONTROLLINGAREA = $user.controllingarea
                    or
                      COSTCENTER = $user.costcenter
                    or
                      COSTCTRRESPONSIBLEUSER = $user.id
                  )
              )
            or
              (
                  void = $user.FAUTHCNTXT1
                and
                   PLANT = $user.plant
              )
          )
      )
    or
      (
          void = $user.FAUTHCNTXT2
        and
           VALUATIONAREA  = $user.VALUATIONAREA
      )
  );
}