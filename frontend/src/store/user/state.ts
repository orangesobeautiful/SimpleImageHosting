export interface UserStateInterface {
  id: number;
  grade: number;
  dataLoaded: boolean;
}

function state(): UserStateInterface {
  return {
    id: -1,
    grade: -1,
    dataLoaded: false,
  };
}

export default state;
