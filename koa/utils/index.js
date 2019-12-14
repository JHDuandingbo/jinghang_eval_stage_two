const lodash = require("lodash")
module.exports={
	queryForQuestions:function (purpose,partNumber, questionCategory, questionType, questionTheme){
		var q = {};
		if(questionTheme){
			q["themeInfo.questionTheme"]= {$regex:questionTheme, $options:"i"};
		}
		if(!lodash.isUndefined(purpose)){
			q = this.queryByPurpose(q,Number(purpose));
		}
		switch(partNumber){
			case "1":
				q.questionPart = "part1";
				break;
			case "2":
				q.questionPart = "part2";
				break;
			case "3":
				q.questionPart = "part3";
				break;
			case "4":
				if(!lodash.isArray(q.$or)){
					q.$or=[];
				}
				q.$or.push({questionPart:"part2"});
				q.$or.push({questionPart:"part3"});
				break;
		}
		switch(questionCategory){
			case "1":
				q["themeInfo.questionCategory"] = "PERSON";
				break;
			case "2":
				q["themeInfo.questionCategory"] = "EVENT";
				break;
			case "3":
				q["themeInfo.questionCategory"] = "OBJECT";
				break;
			case "4":
				q["themeInfo.questionCategory"] = "PLACE";
				break;
		}
		let seasonTag = this.currentSeasonTag();
		if("1" == questionType){//current Season
			q["themeInfo.examTime"]= seasonTag;
		}else if("2"== questionType){
			q["themeInfo.examTime"]= {$ne:seasonTag};
		}

		return q;
	},

	queryByPurpose: function (q, purpose){
		switch(purpose){
			case 0:
				q.isPractice = 0;
				q.isExam = 0;
				break;
			case 1:
				q.isPractice = 1;
				break;
			case 2:
				q.isPractice = 0;
				q.isExam = 1;
				break;
		}
		return q;
	},
	currentSeasonTag:function(){
		let now = new Date();
		let year = now.getFullYear();
		let month = now.getMonth()+1;

		let seasonTag = "";
		if(month >= 1 && month <= 4){
			seasonTag=`${year}年1-4月`;
		}else if(month >=5 && month <= 8){
			seasonTag=`${year}年5-8月`;
		}else{
			seasonTag=`${year}年9-12月`;
		}
		return seasonTag;

	},
	isValidCategory: function(category){
		const list = ["PERSON", "EVENT", "OBJECT", "PLACE"];
		let valid =  list.indexOf(category);
		return (valid === -1 ? 0 :1 );
	}

}
