export class Comment {
  id: number;
  commenter: string;
  body: string;
  createdTime: Date;
}

export class NewComment {
  postId: number;
  commenter: string = "";
  body: string = "";
}
