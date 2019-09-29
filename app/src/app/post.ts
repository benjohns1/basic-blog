import { Comment } from './comment';

export class Post {
  id: number;
  title: string = "";
  body?: string = "";
  author: string = "";
  createdTime: Date = new Date();
  deleted?: boolean = false;
  comments?: Comment[] = [];
}
