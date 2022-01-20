import { MutationTree } from 'vuex';
import { UserStateInterface } from './state';

const mutation: MutationTree<UserStateInterface> = {
  setID(state, id: number) {
    state.id = id;
  },
  setGrade(state, grade: number) {
    state.grade = grade;
  },
  setDataLoaded(state, dataLoaded: boolean) {
    state.dataLoaded = dataLoaded;
  },
};

export default mutation;
