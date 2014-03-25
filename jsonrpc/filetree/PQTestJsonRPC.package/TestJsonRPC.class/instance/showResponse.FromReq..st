private
showResponse: resp FromReq: arithReq
	|  arg |
	arg := arithReq params at: 1.
	Transcript
		show: arithReq method, '(', arg printString, ')' ;
		show: ' : ', resp;
		cr.
