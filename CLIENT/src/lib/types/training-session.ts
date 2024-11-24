export interface TrainingSession {
	id: number;
	trainingSetId: number;
	p1Score: number;
	p2Score: number;
	p3Score: number;
	p4Score: number;
	p5Score: number;
	fullScore: number;
	p1Time: number;
	p2Time: number;
	p3Time: number;
	p4Time: number;
	p5Time: number;
	fullTime: number;
	finished: boolean;
}
